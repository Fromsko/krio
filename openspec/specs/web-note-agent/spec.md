# Web Note Agent Capability

## Requirements

### Requirement: MCP Tool Registration
The system SHALL register a `save_web_note` tool that can be called by AI models through the MCP protocol.

#### Scenario: Tool registration on startup
- **WHEN** the MCP server starts
- **THEN** the system SHALL register the `save_web_note` tool
- **AND** the tool SHALL expose parameters: url (required), tags (optional), folder (optional)

#### Scenario: Tool discovery
- **WHEN** an AI agent queries available tools
- **THEN** the server SHALL return `save_web_note` in the tool list
- **AND** include tool description and input schema

### Requirement: Web Content Fetching
The system SHALL fetch and extract content from provided URLs.

#### Scenario: Successful content fetch
- **WHEN** a valid URL is provided
- **THEN** the system SHALL fetch the webpage content
- **AND** extract the main content (title, body, metadata)
- **AND** remove ads, navigation, and non-essential elements

#### Scenario: URL validation
- **WHEN** an invalid or malformed URL is provided
- **THEN** the system SHALL return an error
- **AND** include a descriptive error message

#### Scenario: Unreachable URL
- **WHEN** the URL cannot be reached (timeout, 404, etc.)
- **THEN** the system SHALL retry up to 3 times
- **AND** return an error if all retries fail

#### Scenario: Security check
- **WHEN** a URL is provided
- **THEN** the system SHALL validate it's not a private/internal address
- **AND** reject URLs pointing to localhost, 127.0.0.1, or private IP ranges

### Requirement: AI Content Summarization
The system SHALL use AI to analyze and summarize the fetched content.

#### Scenario: Generate summary
- **WHEN** web content is successfully fetched
- **THEN** the system SHALL send the content to Claude API
- **AND** request a structured summary including:
  - Title
  - One-sentence summary
  - Key points (3-7 bullet points)
  - Relevant tags
- **AND** receive the AI-generated summary

#### Scenario: Content too long
- **WHEN** the content exceeds token limits
- **THEN** the system SHALL chunk the content
- **AND** process chunks in order
- **AND** merge the summaries

#### Scenario: AI API failure
- **WHEN** the Claude API call fails
- **THEN** the system SHALL retry up to 2 times
- **AND** return an error if all retries fail

### Requirement: Markdown Note Generation
The system SHALL generate well-formatted Markdown notes from the AI summary.

#### Scenario: Standard note format
- **WHEN** the AI summary is received
- **THEN** the system SHALL generate a Markdown note with:
  ```markdown
  ---
  title: [Page Title]
  source: [URL]
  date: [YYYY-MM-DD]
  tags: [tag1, tag2, tag3]
  ---

  # [Page Title]

  > One-sentence summary

  ## Key Points
  - Point 1
  - Point 2
  - Point 3

  ## Details
  [Additional details if needed]
  ```

#### Scenario: Custom tags
- **WHEN** the user provides custom tags
- **THEN** the system SHALL include them in the note frontmatter
- **AND** merge them with AI-generated tags

### Requirement: Obsidian Storage
The system SHALL save the generated Markdown notes to the Obsidian vault.

#### Scenario: Save to default folder
- **WHEN** no folder is specified
- **THEN** the system SHALL save the note to the default inbox folder
- **AND** generate a filename based on the title (sanitized)
- **AND** append timestamp if filename exists

#### Scenario: Save to custom folder
- **WHEN** a folder path is provided
- **THEN** the system SHALL create the folder if it doesn't exist
- **AND** save the note to that folder
- **AND** return the full file path

#### Scenario: File write failure
- **WHEN** file writing fails (permissions, disk full, etc.)
- **THEN** the system SHALL return a descriptive error
- **AND** include the Obsidian vault path for troubleshooting

### Requirement: Response Formatting
The system SHALL return a clear response to the user after note creation.

#### Scenario: Successful note creation
- **WHEN** the note is successfully created
- **THEN** the system SHALL return a success message
- **AND** include the file path
- **AND** include a preview of the note title

#### Scenario: Partial failure
- **WHEN** note creation partially succeeds (e.g., fetched but failed to save)
- **THEN** the system SHALL return an error
- **AND** include details about which step failed
- **AND** suggest troubleshooting steps
