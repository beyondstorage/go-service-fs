name = "fs"

namespace "storage" {
  implement = ["copier", "dir_lister", "mover"]

  new {
    optional = ["work_dir"]
  }

  op "list_dir" {
    optional = ["dir_func", "enable_link_follow", "file_func"]
  }
  op "read" {
    optional = ["offset", "size"]
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
