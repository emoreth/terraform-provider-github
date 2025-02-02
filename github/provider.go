package github

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_TOKEN", nil),
				Description: descriptions["token"],
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_ORGANIZATION", nil),
				Description: descriptions["organization"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_BASE_URL", ""),
				Description: descriptions["base_url"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["insecure"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"github_branch_protection":        resourceGithubBranchProtection(),
			"github_issue_label":              resourceGithubIssueLabel(),
			"github_membership":               resourceGithubMembership(),
			"github_organization_block":       resourceOrganizationBlock(),
			"github_organization_project":     resourceGithubOrganizationProject(),
			"github_organization_webhook":     resourceGithubOrganizationWebhook(),
			"github_project_column":           resourceGithubProjectColumn(),
			"github_repository_collaborator":  resourceGithubRepositoryCollaborator(),
			"github_repository_deploy_key":    resourceGithubRepositoryDeployKey(),
			"github_repository_project":       resourceGithubRepositoryProject(),
			"github_repository_webhook":       resourceGithubRepositoryWebhook(),
			"github_repository":               resourceGithubRepository(),
			"github_team_membership":          resourceGithubTeamMembership(),
			"github_team_repository":          resourceGithubTeamRepository(),
			"github_team":                     resourceGithubTeam(),
			"github_user_gpg_key":             resourceGithubUserGpgKey(),
			"github_user_invitation_accepter": resourceGithubUserInvitationAccepter(),
			"github_user_ssh_key":             resourceGithubUserSshKey(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"github_collaborators": dataSourceGithubCollaborators(),
			"github_ip_ranges":     dataSourceGithubIpRanges(),
			"github_repositories":  dataSourceGithubRepositories(),
			"github_repository":    dataSourceGithubRepository(),
			"github_team":          dataSourceGithubTeam(),
			"github_user":          dataSourceGithubUser(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token": "The OAuth token used to connect to GitHub.",

		"organization": "The GitHub organization name to manage.",

		"base_url": "The GitHub Base API URL",

		"insecure": "Whether server should be accessed " +
			"without verifying the TLS certificate.",
	}
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			Token:        d.Get("token").(string),
			Organization: d.Get("organization").(string),
			BaseURL:      d.Get("base_url").(string),
			Insecure:     d.Get("insecure").(bool),
		}

		meta, err := config.Client()
		if err != nil {
			return nil, err
		}

		meta.(*Organization).StopContext = p.StopContext()

		return meta, nil
	}
}
