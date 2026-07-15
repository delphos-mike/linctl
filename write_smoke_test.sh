#!/bin/bash

# Opt-in WRITE smoke test for linctl.
#
# Unlike smoke_test.sh (read-only), this script exercises create/update/delete
# paths for documents, projects, initiatives, and project status updates. It
# creates throwaway records and cleans them up (documents are trashed via
# `document delete`; projects and initiatives are archived via the `graphql`
# passthrough, since there is no destructive top-level command for them).
#
# It is GATED so it never runs by accident:
#   LINCTL_WRITE_TESTS=1 ./write_smoke_test.sh
#
# Optional:
#   LINCTL_TEST_TEAM=ENG   # team key for throwaway projects (default: first team)

set -u

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

if [ "${LINCTL_WRITE_TESTS:-}" != "1" ]; then
    echo -e "${YELLOW}Skipping write smoke test.${NC}"
    echo "This test creates and deletes real Linear records. To run it:"
    echo "  LINCTL_WRITE_TESTS=1 ./write_smoke_test.sh"
    exit 0
fi

# Binary/command used to invoke linctl. Override with a prebuilt binary via
# LINCTL_BIN=./linctl to avoid repeated `go run` compilation.
LINCTL="${LINCTL_BIN:-go run main.go}"
TS=$(date +%s)
FAILED=0

# Records to clean up on exit (graphql archive mutations).
CLEANUP_PROJECTS=()
CLEANUP_INITIATIVES=()

cleanup() {
    echo -e "\n${YELLOW}Cleaning up throwaway records...${NC}"
    for pid in "${CLEANUP_PROJECTS[@]:-}"; do
        [ -z "$pid" ] && continue
        $LINCTL graphql -q 'mutation($id:String!){ projectArchive(id:$id){ success } }' --var id="$pid" >/dev/null 2>&1 \
            && echo "  archived project $pid" || echo -e "  ${RED}failed to archive project $pid${NC}"
    done
    for iid in "${CLEANUP_INITIATIVES[@]:-}"; do
        [ -z "$iid" ] && continue
        $LINCTL graphql -q 'mutation($id:String!){ initiativeArchive(id:$id){ success } }' --var id="$iid" >/dev/null 2>&1 \
            && echo "  archived initiative $iid" || echo -e "  ${RED}failed to archive initiative $iid${NC}"
    done
}
trap cleanup EXIT

json_field() {
    python3 -c "import sys,json; print(json.load(sys.stdin).get('$1',''))"
}

pass() { echo -e "  ${GREEN}PASS${NC} - $1"; }
fail() { echo -e "  ${RED}FAIL${NC} - $1"; FAILED=$((FAILED + 1)); }

echo "🚀 linctl WRITE smoke test"
echo "================================"

# Resolve a team key for project creation.
TEAM_KEY="${LINCTL_TEST_TEAM:-}"
if [ -z "$TEAM_KEY" ]; then
    TEAM_KEY=$($LINCTL team list 2>/dev/null | awk 'NR>1 {print $1}' | head -1)
fi
if [ -z "$TEAM_KEY" ]; then
    echo -e "${RED}Could not resolve a team key. Set LINCTL_TEST_TEAM.${NC}"
    exit 1
fi
echo "Using team: $TEAM_KEY"

# ---- Documents ----
# Linear requires a document to have a container (project/initiative/team), so
# the throwaway doc is attached to the resolved team.
echo -e "\n${YELLOW}Documents${NC}"
DOC_ID=$($LINCTL document create --title "zzz-write-test-doc-$TS" --team "$TEAM_KEY" --body "initial body" -j 2>/dev/null | json_field id)
if [ -n "$DOC_ID" ]; then
    pass "document create ($DOC_ID)"
else
    fail "document create"
fi

if [ -n "$DOC_ID" ]; then
    TITLE=$($LINCTL document get "$DOC_ID" -j 2>/dev/null | json_field title)
    [ "$TITLE" = "zzz-write-test-doc-$TS" ] && pass "document get" || fail "document get (title=$TITLE)"

    echo "updated body via stdin" | $LINCTL document update "$DOC_ID" >/dev/null 2>&1 \
        && pass "document update (stdin body)" || fail "document update"

    $LINCTL document delete "$DOC_ID" >/dev/null 2>&1 \
        && pass "document delete" || fail "document delete"
fi

# ---- Projects ----
echo -e "\n${YELLOW}Projects${NC}"
PROJ_ID=$($LINCTL project create --name "zzz-write-test-proj-$TS" --team "$TEAM_KEY" --description "throwaway" -j 2>/dev/null | json_field id)
if [ -n "$PROJ_ID" ]; then
    CLEANUP_PROJECTS+=("$PROJ_ID")
    pass "project create ($PROJ_ID)"
else
    fail "project create"
fi

if [ -n "$PROJ_ID" ]; then
    NEWNAME=$($LINCTL project update "$PROJ_ID" --name "zzz-write-test-proj-renamed-$TS" -j 2>/dev/null | json_field name)
    [ "$NEWNAME" = "zzz-write-test-proj-renamed-$TS" ] && pass "project update" || fail "project update (name=$NEWNAME)"

    HEALTH=$($LINCTL project status-update create "$PROJ_ID" --health onTrack --body "smoke update" -j 2>/dev/null | json_field health)
    [ "$HEALTH" = "onTrack" ] && pass "project status-update create" || fail "project status-update create (health=$HEALTH)"
fi

# ---- Initiatives ----
echo -e "\n${YELLOW}Initiatives${NC}"
INIT_ID=$($LINCTL initiative create --name "zzz-write-test-init-$TS" --status Planned -j 2>/dev/null | json_field id)
if [ -n "$INIT_ID" ]; then
    CLEANUP_INITIATIVES+=("$INIT_ID")
    pass "initiative create ($INIT_ID)"
else
    fail "initiative create"
fi

if [ -n "$INIT_ID" ]; then
    STATUS=$($LINCTL initiative update "$INIT_ID" --status Active -j 2>/dev/null | json_field status)
    [ "$STATUS" = "Active" ] && pass "initiative update" || fail "initiative update (status=$STATUS)"

    if [ -n "$PROJ_ID" ]; then
        $LINCTL initiative add-project "$INIT_ID" --project "$PROJ_ID" >/dev/null 2>&1 \
            && pass "initiative add-project" || fail "initiative add-project"

        ATTACHED=$($LINCTL initiative get "$INIT_ID" -j 2>/dev/null | python3 -c "import sys,json; d=json.load(sys.stdin); print(any(p['id']=='$PROJ_ID' for p in (d.get('projects') or {}).get('nodes',[])))")
        [ "$ATTACHED" = "True" ] && pass "initiative get (project attached)" || fail "initiative get (attached=$ATTACHED)"

        $LINCTL initiative remove-project "$INIT_ID" --project "$PROJ_ID" >/dev/null 2>&1 \
            && pass "initiative remove-project" || fail "initiative remove-project"
    fi
fi

# ---- Summary ----
echo -e "\n================================"
if [ "$FAILED" -eq 0 ]; then
    echo -e "${GREEN}✅ All write smoke tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ $FAILED write smoke test(s) failed.${NC}"
    exit 1
fi
