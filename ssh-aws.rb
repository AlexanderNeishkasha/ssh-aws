class SshAws < Formula
  desc "Tool for connection to php servers by private ip"
  homepage "https://github.com/AlexanderNeishkasha/ssh-aws"
  url "https://github.com/AlexanderNeishkasha/ssh-aws/archive/v0.2.2.tar.gz"
  sha256 "158b3b7ae3fe24f31e84d9bfb7e6b2ed3e5af57948dc0b6da752bb421f3ba1a9"

  depends_on "go" => :build

  def install
    bin.install 'ssh-aws'
  end

  test do
    assert_predicate testpath/ssh-aws, :exist?
  end
end
