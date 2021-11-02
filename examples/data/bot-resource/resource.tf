terraform {
  required_providers {
    awslex = {
      version = "0.1"
      source  = "scg.com/va/awslex"
    }
  }
}

data "awslex_bot_resource" "socal_gas_qnabot" {

  id = "EKYRQZCTNM"

  alias = "version7"

}

output "bot" {
  value = "${data.awslex_bot_resource.socal_gas_qnabot}"
}
