terraform {
  required_providers {
    awslex = {
      version = "0.2.0-beta3"
      source  = "localhost/va/awslex"
    }
  }
}

provider "awslex" {
  region = "us-west-2"
}

data "awslex_bot_resource" "socal_gas_qnabot" {

  # id and alias for stable dev
  id = "C5H22UIPWC"
  alias = "latest"

}

output "bot" {
  value = data.awslex_bot_resource.socal_gas_qnabot
}
