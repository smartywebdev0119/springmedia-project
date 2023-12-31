package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

func resourceSourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceLocationCreate,
		ReadContext:   resourceSourceLocationRead,
		UpdateContext: resourceSourceLocationUpdate,
		DeleteContext: resourceSourceLocationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"access_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// may require s3:GetObject
						"access_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"S3_SIGV4", "SECRETS_MANAGER_ACCESS_TOKEN"}, false),
						},
						// SMATC is short for Secrets Manager Access Token Configuration
						"smatc_header_name":       &optionalString,
						"smatc_secret_arn":        &optionalString,
						"smatc_secret_string_key": &optionalString,
					},
				},
			},
			"arn":           &computedString,
			"creation_time": &computedString,
			"default_segment_delivery_configuration_url": &optionalString,
			"http_configuration_url":                     &requiredString,
			"last_modified_time":                         &computedString,
			"segment_delivery_configurations": createOptionalList(
				map[string]*schema.Schema{
					"base_url": &optionalString,
					"name":     &optionalString,
				},
			),
			"name": &requiredString,
			"tags": &optionalTags,
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
		),
	}
}

func resourceSourceLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	var params = getCreateSourceLocationInput(d)

	sourceLocation, err := client.CreateSourceLocation(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the source location: %v", err))
	}
	d.SetId(aws.StringValue(sourceLocation.Arn))

	return resourceSourceLocationRead(ctx, d, meta)
}

func resourceSourceLocationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	resourceName := d.Get("name").(string)
	if len(resourceName) == 0 && len(d.Id()) > 0 {
		resourceArn, err := arn.Parse(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing the name from resource arn: %v", err))
		}
		arnSections := strings.Split(resourceArn.Resource, "/")
		resourceName = arnSections[len(arnSections)-1]
	}
	res, err := client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: aws.String(resourceName)})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the source location: %v", err))
	}

	if err = setSourceLocation(res, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSourceLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")

		resourceName := d.Get("name").(string)
		res, err := client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}

		if err := updateTags(client, res.Arn, oldValue, newValue); err != nil {
			return diag.FromErr(err)
		}
	}

	var params = getUpdateSourceLocationInput(d)
	sourceLocation, err := client.UpdateSourceLocation(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the source location: %v", err))
	}
	d.SetId(aws.StringValue(sourceLocation.Arn))

	return resourceSourceLocationRead(ctx, d, meta)
}

func resourceSourceLocationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)
	sourceLocationName := aws.String(d.Get("name").(string))

	if err := deleteVodSources(sourceLocationName, client); err != nil {
		return diag.FromErr(err)
	}
	if err := deleteLiveSources(sourceLocationName, client); err != nil {
		return diag.FromErr(err)
	}

	_, err := client.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: sourceLocationName})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
	}

	return nil
}
