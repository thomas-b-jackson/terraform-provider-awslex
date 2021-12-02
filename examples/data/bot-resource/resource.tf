terraform {
  required_providers {
    awslex = {
      version = "0.1"
      source  = "thomas-b-jackson/va/awslex"
    }
  }
}

provider "awslex" {
  region = "us-west-2"
}

data "awslex_bot_resource" "socal_gas_qnabot" {

  id = "QU1ORIZZTP"

  alias = "version7"

}

output "bot" {
  value = data.awslex_bot_resource.socal_gas_qnabot
}
