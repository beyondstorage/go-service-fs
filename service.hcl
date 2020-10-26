name = "fs"

namespace "storage" {
  implement = ["copier", "dir_lister", "mover"]

  new {
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["continuation_token", "enable_link_follow"]
  }
  op "read" {
    optional = ["offset", "read_callback_func", "size"]
  }
  op "write" {
    optional = ["size"]
  }
}

pairs {

  pair "enable_link_follow" {
    type = "bool"
  }
}
