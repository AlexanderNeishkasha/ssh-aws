class SshAws < Formula
  desc "Tool for connection to php servers by private ip"
  homepage "https://github.com/AlexanderNeishkasha/ssh-aws"
  url "https://github.com/AlexanderNeishkasha/ssh-aws/archive/v0.1.0.tar.gz"
  sha256 "e59bc27fc385583a119756be674aa320d4d4c4146ed60b474c0e1ee084585c54"

  depends_on "go" => :build

  def install
    bin.install 'ssh-aws'
  end

  test do
    assert_predicate testpath/ssh-aws, :exist?
  end
end
