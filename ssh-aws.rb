class SshAws < Formula
  desc "Tool for connection to php servers by private ip"
  homepage "https://github.com/AlexanderNeishkasha/ssh-aws"
  url "https://github.com/AlexanderNeishkasha/ssh-aws/archive/v0.2.1.tar.gz"
  sha256 "08cb31d398aae881eebd4217aaea9aa09c701d4f167f6d345fccd28c34ffcbbd"

  depends_on "go" => :build

  def install
    bin.install 'ssh-aws'
  end

  test do
    assert_predicate testpath/ssh-aws, :exist?
  end
end
