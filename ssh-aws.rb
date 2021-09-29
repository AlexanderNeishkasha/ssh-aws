class SshAws < Formula
  desc "Tool for connection to php servers by private ip"
  homepage "https://github.com/AlexanderNeishkasha/ssh-aws"
  url "https://github.com/AlexanderNeishkasha/ssh-aws/archive/v0.3.0.tar.gz" :using => GitHubPrivateRepositoryReleaseDownloadStrategy
  sha256 "613a9f828c2113edf8d21e65d0fdfb7d4e785c520f4183de42d008d9738a9840"

  depends_on "go" => :build

  def install
    bin.install 'ssh-aws'
  end

  test do
    assert_predicate testpath/ssh-aws, :exist?
  end
end
