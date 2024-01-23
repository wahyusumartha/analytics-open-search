locals {
  receive_event_function_name = "receive-event"
  receive_event_function_binary_name = local.receive_event_function_name

  receive_event_function = {
    name = local.receive_event_function_name
    src_path = "${path.module}/lambda/${local.receive_event_function_name}"
    binary_name = local.receive_event_function_binary_name
    binary_path = "${path.module}/lambda/${local.receive_event_function_name}/binary/${local.receive_event_function_binary_name}"
    archive_path = "${path.module}/lambda/${local.receive_event_function_name}/binary/${local.receive_event_function_name}.zip"
  }

  ingest_event_function_name = "ingest-event"
  ingest_event_function_binary_name = local.ingest_event_function_name

  ingest_event_function = {
    name = local.ingest_event_function_name
    src_path = "${path.module}/lambda/${local.ingest_event_function_name}"
    binary_name = local.ingest_event_function_binary_name
    binary_path = "${path.module}/lambda/${local.ingest_event_function_name}/binary/${local.ingest_event_function_binary_name}"
    archive_path = "${path.module}/lambda/${local.ingest_event_function_name}/binary/${local.ingest_event_function_name}.zip"
  }
}

output "received_event_function_binary_path" {
  value = local.receive_event_function.binary_path
}

output "ingest_event_function_binary_path" {
  value = local.ingest_event_function.binary_path
}
