
command "terraform" {
  exec = ["/usr/bin/terraform"]

  inherit_env = true
  env = {
    FOO = "bar"
  }

  on_update = ignore
  on_error  = terminate
}
