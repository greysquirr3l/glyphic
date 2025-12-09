# MCP Server Configuration Guide

This guide provides detailed configuration instructions for all Model Context Protocol (MCP) servers
available in the go-moss template. MCP servers extend VS Code's agent mode capabilities by providing
specialized tools for databases, APIs, and development workflows.

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Developer Tools](#developer-tools)
- [Productivity Tools](#productivity-tools)
- [Data & Analytics](#data--analytics)
- [Business Services](#business-services)
- [Cloud & Infrastructure](#cloud--infrastructure)
- [Troubleshooting](#troubleshooting)

## Getting Started

### Prerequisites

- VS Code Insiders with MCP support
- Node.js and npm (for npx commands)
- Python with uvx (for Python-based servers)
- Docker (for containerized servers)

### Basic Configuration

1. Open `.vscode/mcp.json` in your project
2. Uncomment the desired server configuration
3. Replace placeholder values with your actual credentials
4. Restart VS Code to apply changes

### Git Workflow Integration

For enhanced development workflow, consider setting up Git enhancements alongside MCP:

```bash
# Setup comprehensive Git workflow enhancements
./scripts/setup-git-enhancements.sh
```

This complements MCP capabilities with:

- **Conventional Commits**: Consistent commit message format
- **Real-time Validation**: Git hooks with immediate feedback
- **Interactive Commits**: Guided commit creation process
- **CI/CD Integration**: Automated validation in pull requests
- **Team Collaboration**: Standardized workflow for all developers

The Git enhancements work seamlessly with MCP servers, especially GitHub integration,
to provide a complete development environment.

### Configuration Format

```json
{
  "servers": {
    "server-name": {
      "url": "https://api.example.com/mcp",
      "command": "npx",
      "args": ["-y", "package-name"],
      "env": {
        "API_KEY": "your-api-key-here"
      },
      "headers": {
        "Authorization": "Bearer your-token"
      }
    }
  }
}
```

---

## Developer Tools

### GitHub

**Description**: Access GitHub repositories, issues, and pull requests through secure API integration.

**Configuration**:

```json
"github": {
  "url": "https://api.githubcopilot.com/mcp/"
}
```

**Setup**:

1. Ensure you have GitHub Copilot access
2. No additional API keys required (uses Copilot authentication)
3. Provides access to repository data, issues, PRs, and code search

**Capabilities**:

- Repository browsing and file access
- Issue creation and management
- Pull request operations
- Code search across repositories

---

### Figma

**Description**: Extract UI content and generate code from Figma designs.

**Configuration**:

```json
"figma": {
  "url": "http://127.0.0.1:3845/mcp"
}
```

**Setup**:

1. Update Figma desktop app to latest version
2. Enable Dev Mode in your Figma account
3. Start Figma MCP server (automatically runs on port 3845)
4. Ensure Figma desktop app is running

**Capabilities**:

- Design-to-code generation
- Asset extraction from Figma files
- Component analysis and documentation
- Design system integration

---

### Playwright

**Description**: Automate web browsers using accessibility trees for testing and data extraction.

**Configuration**:

```json
"playwright": {
  "command": "npx",
  "args": ["@playwright/mcp@latest"]
}
```

**Setup**:

1. No additional configuration required
2. Installs Playwright automatically via npx
3. Browser binaries downloaded on first use

**Capabilities**:

- Web page automation and testing
- Screenshot and PDF generation
- Form filling and interaction simulation
- Accessibility tree analysis

---

### Sentry

**Description**: Retrieve and analyze application errors and performance issues from Sentry projects.

**Configuration**:

```json
"sentry": {
  "url": "https://mcp.sentry.dev/mcp"
}
```

**Setup**:

1. Create Sentry account at [sentry.io](https://sentry.io)
2. Generate API token in Sentry Settings → API
3. Configure authentication through Sentry's MCP interface

**Capabilities**:

- Error tracking and analysis
- Performance monitoring data
- Release and deployment tracking
- Issue assignment and resolution

---

### Hugging Face

**Description**: Access models, datasets, and Spaces on the Hugging Face Hub.

**Configuration**:

```json
"huggingface": {
  "url": "https://hf.co/mcp"
}
```

**Setup**:

1. Create account at [huggingface.co](https://huggingface.co)
2. Generate access token in Settings → Access Tokens
3. Authentication handled through Hugging Face's MCP service

**Capabilities**:

- Model discovery and analysis
- Dataset exploration and download
- Spaces interaction and deployment
- Model inference and evaluation

---

### Additional Developer Tools

**DeepWiki**: Query GitHub repositories indexed on DeepWiki

**MarkItDown**: Convert files (PDF, Word, Excel, images, audio) to Markdown

**Microsoft Docs**: Search Microsoft Learn and Azure documentation

**Context7**: Get up-to-date library and framework documentation

**ImageSorcery**: Local image processing with computer vision

**Codacy**: Code quality and security analysis with SAST scanning

---

## Productivity Tools

### Notion

**Description**: View, search, create, and update Notion pages and databases.

**Configuration**:

```json
"notion": {
  "command": "npx",
  "args": ["-y", "@notionhq/notion-mcp-server"],
  "env": {
    "NOTION_TOKEN": "your-notion-integration-token"
  }
}
```

**Setup**:

1. Create Notion integration at [notion.so/integrations](https://www.notion.so/profile/integrations)
2. Copy the Internal Integration Token
3. Share relevant pages/databases with your integration

**Capabilities**:

- Page creation and editing
- Database queries and updates
- Content search across workspace
- Block-level content manipulation

---

### Linear

**Description**: Create, update, and track issues in Linear's project management platform.

**Configuration**:

```json
"linear": {
  "url": "https://mcp.linear.app/sse"
}
```

**Setup**:

1. Create Linear account at [linear.app](https://linear.app)
2. Authentication handled through Linear's OAuth flow
3. Grant necessary permissions for your workspace

**Capabilities**:

- Issue creation and management
- Project tracking and updates
- Team collaboration features
- Workflow automation

---

### Additional Productivity Tools

**Sequential Thinking**: Break down complex tasks into manageable steps

**Memory**: Store and retrieve contextual information across sessions

**Asana**: Create and manage tasks, projects, and comments

**Atlassian**: Connect to Jira and Confluence for issue tracking

**Zapier**: Create workflows across 30,000+ connected apps

**Monday.com**: Project management with boards, items, and teams

---

## Data & Analytics

### DuckDB

**Description**: Query and analyze data in DuckDB databases locally and in the cloud.

**Configuration**:

```json
"duckdb": {
  "command": "uvx",
  "args": ["mcp-server-duckdb", "--db-path", "/path/to/your/database.duckdb"]
}
```

**Setup**:

1. Install Python and uvx
2. Create or specify existing DuckDB database file
3. Ensure read/write permissions for database file

**Capabilities**:

- SQL query execution
- Data analysis and aggregation
- CSV/Parquet file import
- Analytics and reporting

---

### PostHog

**Description**: Access PostHog analytics to create annotations and retrieve product usage insights.

**Configuration**:

```json
"posthog": {
  "url": "https://mcp.posthog.com/sse",
  "headers": {
    "Authorization": "Bearer your-posthog-api-key"
  }
}
```

**Setup**:

1. Create PostHog account at [posthog.com](https://posthog.com)
2. Generate API key with `annotation:write` and `project:read` scopes
3. Note your project ID for queries

**Capabilities**:

- Product analytics data access
- Event tracking and analysis
- User behavior insights
- Custom annotation creation

---

### Additional Data & Analytics Tools

**Neon**: Manage and query Neon Postgres databases with natural language

**Apify**: Extract data from websites and automate workflows

**Microsoft Clarity**: Access analytics data including heatmaps and session recordings

**Firecrawl**: Advanced web scraping, crawling, and structured data extraction

**Prisma Postgres**: Database operations with Prisma ORM and PostgreSQL

**MongoDB**: Database operations including queries, collections, and aggregation

---

## Business Services

### Stripe

**Description**: Create customers, manage subscriptions, and generate payment links through Stripe APIs.

**Configuration**:

```json
"stripe": {
  "url": "https://mcp.stripe.com/",
  "headers": {
    "Authorization": "Bearer your-stripe-secret-key"
  }
}
```

**Setup**:

1. Create Stripe account at [stripe.com](https://stripe.com)
2. Get API keys from [dashboard.stripe.com/apikeys](https://dashboard.stripe.com/apikeys)
3. Use test keys for development, live keys for production

**Capabilities**:

- Customer management
- Payment processing
- Subscription handling
- Invoice generation

---

### Additional Business Services

**PayPal**: Create invoices, process payments, and access transaction data

**Square**: Process payments and manage customers through Square's API

**Intercom**: Access customer conversations and support tickets

**Wix**: Build and manage Wix sites with eCommerce and payment features

**Webflow**: Create and manage websites, collections, and content

---

## Cloud & Infrastructure

### Azure

**Description**: Manage Azure resources, query databases, and access Azure services.

**Configuration**:

```json
"azure": {
  "command": "npx",
  "args": ["-y", "@azure/mcp@latest", "server", "start"]
}
```

**Setup**:

1. Have Azure subscription
2. Install Azure CLI and authenticate
3. Configure service principal or user authentication

**Capabilities**:

- Resource management
- Database operations
- Service configuration
- Monitoring and analytics

---

### Additional Cloud & Infrastructure Tools

**Convex**: Access Convex backend databases and functions for real-time operations

**Azure DevOps**: Manage projects, work items, repositories, builds, and releases

**Terraform**: Infrastructure as Code management with plan, apply, destroy operations

---

## Troubleshooting

### Common Issues

#### Server Won't Start

- **Check prerequisites**: Ensure Node.js, Python, or Docker are installed as required
- **Verify credentials**: Double-check API keys and tokens
- **Check network**: Ensure internet connectivity for cloud-based servers
- **Review logs**: Check VS Code developer console for error messages

#### Authentication Failures

- **Regenerate tokens**: Try creating new API keys or tokens
- **Check permissions**: Ensure tokens have required scopes
- **Verify expiration**: Check if tokens have expired
- **Review configuration**: Ensure environment variables are set correctly

#### Command Not Found Errors

- **Install dependencies**: Run `npm install -g` for global packages
- **Check PATH**: Ensure command-line tools are in your PATH
- **Use full paths**: Specify complete paths to executables if needed

### Performance Optimization

#### Reduce Startup Time

- **Cache dependencies**: Use package managers' cache features
- **Minimize servers**: Only enable servers you actively use
- **Use local servers**: Prefer local over remote servers when possible

#### Memory Management

- **Monitor usage**: Check VS Code's memory consumption
- **Restart periodically**: Restart servers that consume too much memory
- **Optimize queries**: Use efficient queries for database servers

### Security Best Practices

#### Credential Management

- **Use environment variables**: Store sensitive data in environment variables
- **Rotate keys regularly**: Change API keys and tokens periodically
- **Limit permissions**: Use least-privilege principle for tokens
- **Secure storage**: Use secure credential storage systems

#### Network Security

- **Use HTTPS**: Ensure all connections use secure protocols
- **Validate certificates**: Don't skip certificate validation
- **Monitor access**: Keep track of API usage and access patterns

### Getting Help

#### Official Resources

- **VS Code MCP Documentation**: [aka.ms/vscode-website-mcp-docs](https://aka.ms/vscode-website-mcp-docs)
- **Model Context Protocol**: [modelcontextprotocol.io](https://modelcontextprotocol.io)
- **Individual Server Documentation**: Check each server's GitHub repository

#### Community Support

- **VS Code GitHub**: Report issues and feature requests
- **Discord/Slack**: Join community channels for real-time help
- **Stack Overflow**: Search for existing solutions and ask questions

---

## Conclusion

This configuration guide provides comprehensive setup instructions for all available MCP servers.
Start with servers that match your immediate needs, and gradually expand your toolkit as your
projects grow in complexity.

Remember to:

- Keep credentials secure and up to date
- Monitor server performance and resource usage
- Stay updated with server versions and new features
- Contribute back to the community with improvements and feedback

For the latest updates and new MCP servers, check the
[official VS Code MCP documentation](https://aka.ms/vscode-website-mcp-docs).
