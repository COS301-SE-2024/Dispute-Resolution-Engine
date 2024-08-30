# Coding Standards

## Formatting Tools and Configurations

### Prettier
We use Prettier for code formatting to ensure consistency across different files and contributors for all our Javascript/Typescript project. Our prettier configuartion is as follows:

```json
{
  "tabWidth": 2,
  "useTabs": false,
  "printWidth": 100,
  "singleQuote": false,
  "bracketSameLine": false,
  "arrowParens": "always",
  "endOfLine": "lf"
}
```

### Go Format
For our Go-based projects, we use `gofmt` along with the opinionated formatting of the Go
language, which ensures that all our Go code follow the style prescribed by the Go development
team.

## Git Conventions

### Branch Naming Rules
Branches in our repository are divided into two types of branches:
- `feat/xxx`: Features to be added to the code base
- `docs/xxx`: Documentation-based branches
- `hotfix/xxx`: Urgent changes to be made to main 

### Commit Naming Rules
Commit messages should be concise and descriptive. The format is:
- `feat: xxx` for additions that add features to the project
- `refac: xxx` for changes made to the code with no additional features added
- `fix: xxx` for bug fixes
- `chore: xxx` for general structural organization without features
- `docs: xxx` for changes made to documentation in general

## File Structure
Our repository follows a monorepo structure, with each part of the project defined
in a sub-folder in the root of the project:

```
├── api       <---- Go API
├── frontend  <---- Next.js frontend
└── initdb    <---- SQL scripts to initialize database
```

## Conclusion
Adhering to these coding standards and conventions will help ensure that our codebase remains clean, consistent, and maintainable. These guidelines should be reviewed and updated as necessary to accommodate new tools or practices.
