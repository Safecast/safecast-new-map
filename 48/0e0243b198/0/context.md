# Session Context

## User Prompts

### Prompt 1

Now we use Cloudfront for DNS of simplemap.safecast.org, does rsync still works with ssh to simplemap.safecast.org?

### Prompt 2

Can you update all the documentaion in the repo that need to be updated for this chnage?

### Prompt 3

Can you do the same for the repo I add for the MCP server? And check if we need to modify the MCP server to work correctly with the new Cloudfront setup?

### Prompt 4

Can you do this?:

You must update the MAP_SERVER_HOST secret in the MCP server repository:

Go to safecast-map-MCP repository on GitHub
Settings → Secrets and variables → Actions
Update MAP_SERVER_HOST from simplemap.safecast.org to 65.108.24.131

### Prompt 5

Is there a way to get CloudFront  to route the ssh trafic to the server with the domian name?

