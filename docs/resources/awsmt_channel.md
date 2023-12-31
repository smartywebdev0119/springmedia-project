# Resource: awsmt_channel

Provides an Elemental MediaTailor Channel.

## Example Usage

```terraform
resource "awsmt_channel" "example" {
  name = "example-channel"
  outputs {
    manifest_name                = "default"
    source_group                 = "default"
    hls_manifest_windows_seconds = 30
  }
  playback_mode = "LOOP"
  tier          = "BASIC"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the channel.
- `channel_state` - (Optional) The state of the channel. Can be either `RUNNING` or `STOPPED`.
- `filler_slate` – (Optional) The slate used to fill gaps between programs in the schedule. You must configure filler slate if your channel uses the LINEAR PlaybackMode.
  - `source_location_name` - (Optional) The name of the source location where the slate VOD source is stored.
  - `vod_source_name` - (Optional) The slate VOD source name. The VOD source must already exist in a source location before it can be used for slate.
- `outputs` – (Optional) The channel's output properties.
  - `dash_manifest_windows_seconds` - (Optional) The total duration (in seconds) of each dash manifest.
  - `dash_min_buffer_time_seconds` - (Optional) Minimum amount of content (measured in seconds) that a player must keep available in the buffer.
  - `dash_min_update_period_seconds` - (Optional) Minimum amount of time (in seconds) that the player should wait before requesting updates to the manifest.
  - `dash_suggested_presentation_delay_seconds` - (Optional) Amount of time (in seconds) that the player should be from the live point at the end of the manifest.
  - `hls_manifest_windows_seconds` - (Optional) The total duration (in seconds) of each hls manifest.
  - `manifest_name` - (Required) The name of the manifest for the channel. The name appears in the PlaybackUrl.
- `playback_mode` - (Required) The type of playback mode for this channel. Can be either LINEAR or LOOP.
- `policy` - (Required) The IAM policy for the channel.
- `source_group` - (Required) A string used to match which HttpPackageConfiguration is used for each VodSource.
- `tags` - (Optional) Key-value mapping of resource tags.
- `tier` - (Required) The tier for this channel. STANDARD tier channels can contain live programs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `arn` - The ARN of the channel.
- `creation_time` - The timestamp of when the channel was created.
- `last_modified_time` - The timestamp of when the channel was last modified.
- `outputs` – The channel's output properties.
  - `playback_url` - The URL used for playback by content players.

## Import

Channels can be imported using their ARN as identifier. For example:

```shell
  $ terraform import awsmt_channel.example arn:aws:mediatailor:us-east-1:000000000000:channel/example
```
