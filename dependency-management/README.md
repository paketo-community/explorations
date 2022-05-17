# Dependency Management Exploration
This exploration is primarily to investigate solutions and architecture options
to put an RFC together to address
https://github.com/paketo-community/explorations/issues/8.


## Ideal outcomes
1.  Design a system in which language-family experts are the decision
    makers/maintainers of language-specific concerns:
    - where dependencies come from
    - are they compiled or processed in anyway, or are they taken directly from source?
    - how often dependencies are updated
    - decide which stacks/OSe

2.  Design a system in which Dependencies sub-team members are decision
    makers/maintainers of general dependency-related process and automation:
    - set up infrastructure for any dependencies that need to be hosted by the project
    - set up generic automation to handle dependency updates that can be used in language-family workflows
    - allow for an easy way to get a dependency locally (especially if compiled)
    - establish project-wide guidelines/best practices for our dependency management approach

3. Design a system in which permissions are clearly delegated to the right parties at the right level.
   - if a change is needed to language-specific logic, language maintainers should be the owners of that code
   - if a change is needed to our project-wide approach, dependencies maintaienrs should be the owners
   - infrastructure/buckets are owned by the project, but users can set their own instances up easily


TODOs:
- Who will pay for buckets?
- GHA slowdowns?
- Stacks and architecture questions
- Test out ARM64 on GHA
