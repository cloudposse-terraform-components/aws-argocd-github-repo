locals {
  github_token = local.enabled ? (
    var.use_local_github_credentials ? null : coalesce(var.github_token_override, data.aws_ssm_parameter.github_api_key[0].value)
  ) : ""

  deploy_key_environments = local.deploy_keys_enabled ? local.environments : {}
}

data "aws_ssm_parameter" "github_api_key" {
  count           = local.enabled && !var.use_local_github_credentials ? 1 : 0
  name            = var.ssm_github_api_key
  with_decryption = true
}

module "store_write" {
  source  = "cloudposse/ssm-parameter-store/aws"
  version = "0.11.0"

  parameter_write = [for k, v in local.deploy_key_environments :
    {
      name        = format(var.ssm_github_deploy_key_format, k)
      value       = tls_private_key.default[k].private_key_pem
      type        = "SecureString"
      overwrite   = true
      description = github_repository_deploy_key.default[k].title
    }
  ]

  context = module.this.context
}

provider "github" {
  base_url = var.github_base_url
  owner    = var.github_organization
  token    = local.github_token
}
