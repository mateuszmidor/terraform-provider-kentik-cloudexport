package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kentik/community_sdk_golang/apiv6/kentikapi/cloudexport"
)

// schemaMode determines if we want a schema for:
// -reading single item - we need to provide "id" of the item to read, everything else is provided by server
// -reading list of items - we don't need to provide a thing, everything is provided by server
// -creating new item - we need to provide a bunch of obligatory attributes, the rest is provided by the server
type schemaMode int

const (
	READ_SINGLE schemaMode = iota
	READ_LIST
	CREATE
)

// CloudExportSchema reflects V202101beta1CloudExport type and defines a CloudExport item used in terraform .tf files
// Note: currently, nesting an object is only possible by using single-item List element (Terraform limitation)
func makeCloudExportSchema(mode schemaMode) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: mode == CREATE || mode == READ_LIST, // provided by server on creating/listing items
			Required: mode == READ_SINGLE,                 // provided by user in order to read single item
		},
		"type": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Required:    mode == CREATE,                           // provided by user on create
			Description: "One of [CLOUD_EXPORT_TYPE_KENTIK_MANAGED, CLOUD_EXPORT_TYPE_CUSTOMER_MANAGED]",
		},
		"enabled": &schema.Schema{
			Type:     schema.TypeBool,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Required: mode == CREATE,                           // provided by user on create
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
		},
		"api_root": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true, // always provided by server
		},
		"flow_dest": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true, // always provided by server
		},
		"plan_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Required: mode == CREATE,                           // provided by user on create
		},
		"cloud_provider": &schema.Schema{
			Type:        schema.TypeString,
			Computed:    mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Required:    mode == CREATE,                           // provided by user on create
			Description: "One of [aws, azure, ibm, gce, bgp]",
		},
		"aws": &schema.Schema{
			// nested object
			Type:     schema.TypeList,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"iam_role_arn": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"region": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"delete_after_read": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"multiple_buckets": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
				},
			},
		},
		"azure": &schema.Schema{
			// nested object
			Type:     schema.TypeList,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"location": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"resource_group": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"storage_account": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"subscription_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"security_principal_enabled": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
				},
			},
		},
		"bgp": &schema.Schema{
			// nested object
			Type:     schema.TypeList,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"apply_bgp": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"use_bgp_device_id": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"device_bgp_type": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
				},
			},
		},
		"gce": &schema.Schema{
			// nested object
			Type:     schema.TypeList,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"project": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
					"subscription": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
				},
			},
		},
		"ibm": &schema.Schema{
			// nested object
			Type:     schema.TypeList,
			Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
			Optional: mode == CREATE,                           // optionally provided by user on create
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"bucket": &schema.Schema{
						Type:     schema.TypeString,
						Computed: mode == READ_SINGLE || mode == READ_LIST, // provided by server on read
						Required: mode == CREATE,                           // provided by user on create
					},
				},
			},
		},
		"current_status": &schema.Schema{
			// nested object
			Type:     schema.TypeList,
			Computed: true, // always provided by server
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"status": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"error_message": &schema.Schema{
						Type:     schema.TypeString,
						Computed: true,
					},
					"flow_found": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
					"api_access": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
					"storage_account_access": &schema.Schema{
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
	}
}

// cloudExportToMap is used for API get operation to fill terraform resource from cloudexport item
func cloudExportToMap(e *cloudexport.V202101beta1CloudExport) map[string]interface{} {
	o := make(map[string]interface{})
	if e == nil {
		return o
	}

	o["id"] = e.Id
	o["type"] = e.Type
	o["enabled"] = e.Enabled
	o["name"] = e.Name
	o["description"] = e.Description
	o["api_root"] = e.ApiRoot
	o["flow_dest"] = e.FlowDest
	o["plan_id"] = e.PlanId
	o["cloud_provider"] = e.CloudProvider

	if e.Aws != nil {
		aws := make(map[string]interface{})
		aws["bucket"] = e.Aws.Bucket
		aws["iam_role_arn"] = e.Aws.IamRoleArn
		aws["region"] = e.Aws.Region
		aws["delete_after_read"] = e.Aws.DeleteAfterRead
		aws["multiple_buckets"] = e.Aws.MultipleBuckets
		o["aws"] = []interface{}{aws}
	}

	if e.Azure != nil {
		azure := make(map[string]interface{})
		azure["location"] = e.Azure.Location
		azure["resource_group"] = e.Azure.ResourceGroup
		azure["storage_account"] = e.Azure.StorageAccount
		azure["subscription_id"] = e.Azure.SubscriptionId
		azure["security_principal_enabled"] = e.Azure.SecurityPrincipalEnabled
		o["azure"] = []interface{}{azure}
	}

	if e.Bgp != nil {
		bgp := make(map[string]interface{})
		bgp["apply_bgp"] = e.Bgp.ApplyBgp
		bgp["use_bgp_device_id"] = e.Bgp.UseBgpDeviceId
		bgp["device_bgp_type"] = e.Bgp.DeviceBgpType
		o["bgp"] = []interface{}{bgp}
	}

	if e.Gce != nil {
		gce := make(map[string]interface{})
		gce["project"] = e.Gce.Project
		gce["subscription"] = e.Gce.Subscription
		o["gce"] = []interface{}{gce}
	}

	if e.Ibm != nil {
		ibm := make(map[string]interface{})
		ibm["bucket"] = e.Ibm.Bucket
		o["ibm"] = []interface{}{ibm}
	}

	if e.CurrentStatus != nil {
		current_status := make(map[string]interface{})
		current_status["status"] = e.CurrentStatus.Status
		current_status["error_message"] = e.CurrentStatus.ErrorMessage
		current_status["flow_found"] = e.CurrentStatus.FlowFound
		current_status["api_access"] = e.CurrentStatus.ApiAccess
		current_status["storage_account_access"] = e.CurrentStatus.StorageAccountAccess
		o["current_status"] = []interface{}{current_status}
	}

	return o
}

// resourceDataToCloudExport is used for API create/update operations to fill cloudexport item from terraform resource
func resourceDataToCloudExport(d *schema.ResourceData) (*cloudexport.V202101beta1CloudExport, error) {
	export := cloudexport.NewV202101beta1CloudExport()

	export.SetId(d.Get("id").(string))
	export.SetType(cloudexport.V202101beta1CloudExportType(d.Get("type").(string)))
	export.SetEnabled(d.Get("enabled").(bool))
	export.SetName(d.Get("name").(string))
	export.SetDescription(d.Get("description").(string))
	export.SetApiRoot(d.Get("api_root").(string))
	export.SetFlowDest(d.Get("flow_dest").(string))
	export.SetPlanId(d.Get("plan_id").(string))

	// validation: for any given cloud_provider, there should also be an object of the same name, containing configuration details
	// eg for cloud_provider="ibm", ibm{...} object should be provided
	cloudProvider := d.Get("cloud_provider").(string)
	providerObj, ok := d.GetOk(cloudProvider)
	if !ok {
		return nil, fmt.Errorf("for cloud_provider=%[1]s, there should also be %[1]s{...} attribute provided", cloudProvider)
	}
	export.SetCloudProvider(cloudProvider)
	providerDef := providerObj.([]interface{})[0] // extract nested object under index 0. Terraform clumsyness
	providerMap := providerDef.(map[string]interface{})
	switch cloudProvider {
	case "aws":
		{
			aws := *cloudexport.NewV202101beta1AwsProperties()
			aws.SetBucket(providerMap["bucket"].(string))
			aws.SetIamRoleArn(providerMap["iam_role_arn"].(string))
			aws.SetRegion(providerMap["region"].(string))
			aws.SetDeleteAfterRead(providerMap["delete_after_read"].(bool))
			aws.SetMultipleBuckets(providerMap["multiple_buckets"].(bool))
			export.SetAws(aws)
		}
	case "azure":
		{
			azure := *cloudexport.NewV202101beta1AzureProperties()
			azure.SetLocation(providerMap["location"].(string))
			azure.SetResourceGroup(providerMap["resource_group"].(string))
			azure.SetStorageAccount(providerMap["storage_account"].(string))
			azure.SetSubscriptionId(providerMap["subscription_id"].(string))
			azure.SetSecurityPrincipalEnabled(providerMap["security_principal_enabled"].(bool))
			export.SetAzure(azure)
		}
	case "bgp":
		{
			bgp := *cloudexport.NewV202101beta1BgpProperties()
			bgp.SetApplyBgp(providerMap["apply_bgp"].(bool))
			bgp.SetUseBgpDeviceId(providerMap["use_bgp_device_id"].(string))
			bgp.SetDeviceBgpType(providerMap["device_bgp_type"].(string))
			export.SetBgp(bgp)
		}
	case "gce":
		{
			gce := *cloudexport.NewV202101beta1GceProperties()
			gce.SetProject(providerMap["project"].(string))
			gce.SetSubscription(providerMap["subscription"].(string))
			export.SetGce(gce)
		}
	case "ibm":
		{
			ibm := *cloudexport.NewV202101beta1IbmProperties()
			ibm.SetBucket(providerMap["bucket"].(string))
			export.SetIbm(ibm)
		}
	default:
		return nil, fmt.Errorf("cloud_provider should be one of [aws, azure, ibm, gce, bgp], got: %q", cloudProvider)
	}

	return export, nil
}
