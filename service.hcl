name = "fs"

namespace "storage" {
  implement = ["copier", "dir_lister", "mover"]

  new {
    optional = ["pair_policy", "work_dir"]
  }

  op "list_dir" {
    optional = ["continuation_token"]
  }
  op "read" {
    optional = ["offset", "read_callback_func", "size"]
  }
  op "write" {
    optional = ["offset", "read_callback_func", "size"]
  }
}
