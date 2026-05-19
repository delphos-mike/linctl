class Linctl < Formula
  desc "Comprehensive command-line interface for Linear's API"
  homepage "https://github.com/delphos-mike/linctl"
  url "https://github.com/delphos-mike/linctl/archive/refs/tags/v0.3.0.tar.gz"
  sha256 "3a1dbf83ef9877aa528635b8f5ee376ce01780a7bb94e588dc6845d5354d41a9"
  license "MIT"
  head "https://github.com/delphos-mike/linctl.git", branch: "main"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w -X github.com/delphos-mike/linctl/cmd.version=#{version}")
  end

  test do
    # Test version output
    assert_match "linctl version #{version}", shell_output("#{bin}/linctl --version")
    
    # Test help command
    assert_match "A comprehensive CLI tool for Linear", shell_output("#{bin}/linctl --help")
  end
end