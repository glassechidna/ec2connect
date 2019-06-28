# `ec2connect`

![render1561718563616](https://user-images.githubusercontent.com/369053/60337186-ad908180-99e5-11e9-9db3-de8b0d353739.gif)

In [June 2019][rel-notes], AWS released [EC2 Instance Connect][docs] - a way of
authenticating SSH sessions using AWS IAM policies. This **massively** improves
security by removing the need for sharing SSH private keys. It also improves
reliability by removing the need for any workarounds to avoid sharing keys!

AWS did release an [`mssh`][mssh] tool, but it's not as nice as it could be.
`ec2connect` improves upon it:

* Doesn't require Python to be installed. Single binary available for Mac, Linux
  and Windows.
* Doesn't require a new command to be remembered - just `ssh ec2-user@host` as 
  normal.
* Integrates nicely with every other tool - any tool that relies on SSH (e.g. `git`)
  will work out of the box due to the above.

## Installation

* Mac: `brew install glassechidna/taps/ec2connect`
* Windows: `scoop bucket add glassechidna https://github.com/glassechidna/scoop-bucket.git; scoop install ec2connect`
* Otherwise get the latest build from the [Releases][releases] tab.

## Usage

On first time usage, run `ec2connect setup`. This sets up your SSH configuration
to use `ec2connect` to connect to your instances. You only need to run this once.

Now, connect to your instances using `ssh <user>@<instance id>`. For example:

```
# regular ssh connection
ssh ec2-user@i-000abc124def

# in a different region
AWS_REGION=us-west-2 ssh ec2-user@i-000abc124def

# with a profile
AWS_PROFILE=mycompany ssh ec2-user@i-000abc124def

# with port-forwarding. the possibilities are endless!
ssh -L 2375:127.0.0.1:2375 ec2-user@i-000abc124def
```

[rel-notes]: https://aws.amazon.com/about-aws/whats-new/2019/06/introducing-amazon-ec2-instance-connect/
[docs]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Connect-using-EC2-Instance-Connect.html
[mssh]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-connect-set-up.html#ec2-instance-connect-install-eic-CLI
