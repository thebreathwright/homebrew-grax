class Grax < Formula
  desc "Grax model/runtime contract and Modelfile assets"
  homepage "https://github.com/ollama/ollama"
  license "MIT"
  head "https://github.com/thebreathwright/homebrew-grax.git", branch: "main"

  def install
    pkgshare.install "Modelfile"
    pkgshare.install "docs/grax_package.md"
    pkgshare.install "model/parsers/grax.go"
    pkgshare.install "model/renderers/grax.go"
  end

  def caveats
    <<~EOS
      Grax installs model-contract files only.
      Use them with the brullama CLI and your local model runtime.
    EOS
  end

  test do
    assert_predicate pkgshare/"Modelfile", :exist?
    assert_predicate pkgshare/"docs/grax_package.md", :exist?
  end
end
