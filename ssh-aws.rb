class SshAws < Formula
  desc "Tool for connection to php servers by private ip"
  homepage "https://github.com/AlexanderNeishkasha/ssh-aws"
  url "https://github.com/AlexanderNeishkasha/ssh-aws/archive/v0.2.0.tar.gz"
  sha256 "870f0d3d90910c7d59c21411d19998f675f5dca082f150ec5ab4bc13672fbeeb"

  depends_on "go" => :build

  def install
    bin.install 'ssh-aws'
  end

  test do
    assert_predicate testpath/ssh-aws, :exist?
  end
end
