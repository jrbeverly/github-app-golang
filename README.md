# GitHub App in Golang Prototyping

GitHub App written in Golang with the AWS & GitHub integrations split away, to try and encode the core 'concepts' solely into the lib/ component

## Notes

- Development (& Testing) should support a non-smee way of development
- How would integration/infrastructure/e2e tests work with this kind of system?
- Would it make more sense to have an OpenAPI system, with the GitHub integration performing the interface?
- Is this similar to Hubot, in that having an interface layer/service makes more sense, as it avoids the requirement for direct interface? (This can be within the system itself)
