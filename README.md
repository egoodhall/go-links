# go-links

Service for managing easy-to-remember redirects to commonly used tools. All links are configured in a static file, to keep the operations simple.

### Example Configuration

```yaml
address: :8443
targets:
-
  title: GitHub
  description: Online code repositories
  aliases:
  - git
  - gh
  urls:
  - https://github.com/egoodhall
  - https://github.com/egoodhall/:repo
  - https://github.com/:org/:repo
```

This will start the service with the following paths:
- `GET /git`
- `GET /git/:repo`
- `GET /git/:org/repo`
- `GET /gh`
- `GET /gh/:repo`
- `GET /gh/:org/repo`
