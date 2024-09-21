# YouTube Data Fetcher

This Go project allows you to extract YouTube video information such as transcripts, duration, and comments using the YouTube API. It retrieves the video transcript, parses the video duration, and optionally fetches comments. The project uses the `soup` library for HTML parsing and the official YouTube Data API.

## Features

- Extract video transcripts.
- Fetch video comments.
- Parse and output video duration in minutes.
- Outputs data in JSON format.

## Prerequisites

- Go 1.22 or higher
- YouTube API key (see instructions below to obtain one)

## Installation

### Option 1: Install via `go install`

You can install the project directly using the following command:

```bash
go install github.com/danielmiessler/yt@latest
```

This will install the program to your `$GOPATH/bin`, making it available globally in your terminal.

### Option 2: Manual Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/danielmiessler/yt.git
   cd yt
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Install the program:

   ```bash
   go install
   ```

4. Create a configuration file for your API key:

   ```bash
   mkdir -p ~/.config/fabric
   echo "YOUTUBE_API_KEY=[Your API Key]" >> ~/.config/fabric/.env
   ```

## How to Get a YouTube API Key

To access the YouTube Data API, you need to obtain an API key. Follow these steps:

1. **Log in to Google Developers Console**: Go to the [Google Developers Console](https://console.cloud.google.com/) and sign in using your Google account. If you don’t have one, create it.

2. **Create a New Project**: Once logged in, click the "Create Project" button at the top-right. Name the project and click "Create."

3. **Enable the YouTube Data API**: 
   - Navigate to "APIs & Services" > "Library."
   - In the API Library, search for "YouTube Data API v3" and select it.
   - Click "Enable."

4. **Create Credentials**:
   - Go to "APIs & Services" > "Credentials."
   - Click "Create Credentials" and select "API Key." Google will generate an API key for you.

5. **Copy and Save the API Key**: Copy the generated API key and store it securely, as you will need it to access the YouTube Data API.

6. (Optional) **Restrict API Key Usage**: To prevent unauthorized use of your API key, consider restricting it to specific IP addresses or services like YouTube Data API v3.

For more information, refer to [Google’s official guide](https://console.cloud.google.com/).

## Usage

### Run the application

```bash
go run main.go [YouTube URL] [options]
```

### Options

- `-duration`: Output only the video duration in minutes.
- `-transcript`: Output only the video transcript.
- `-comments`: Output the comments of the video.
- `-lang`: Specify the language for the transcript (default: English).

### Example

Fetch transcript and video information:

```bash
go run main.go "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

Fetch comments:

```bash
go run main.go "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -comments
```

Fetch video duration:

```bash
go run main.go "https://www.youtube.com/watch?v=dQw4w9WgXcQ" -duration
```

## Environment Variables

- `YOUTUBE_API_KEY`: You must provide a valid API key for accessing the YouTube Data API.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
