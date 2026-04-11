+++
title = "Setup"
weight = 10
+++

Use this page for installation steps, platform notes, or entry checks that do not belong in the narrative flow of an episode.

For OS-specific instructions, prefer Hugo-native tabs over large walls of repeated setup text. `hugo-styles` enables Hextra's synced tabs by default, so repeated choices like `macOS`, `Linux`, and `Windows` stay aligned across a page.

## Operating system setup

{{< tabs >}}
{{< tab name="macOS" selected=true >}}
```bash
brew install hugo go
```

Use Homebrew when you want a quick local setup for lesson authoring.
{{< /tab >}}
{{< tab name="Linux" >}}
```bash
sudo apt update
sudo apt install hugo golang
```

Package names vary slightly by distribution, so note that in downstream lessons when needed.
{{< /tab >}}
{{< tab name="Windows" >}}
```powershell
winget install Hugo.Hugo.Extended
winget install GoLang.Go
```

If your learners use corporate-managed machines, add a short note about alternate installation paths.
{{< /tab >}}
{{< /tabs >}}

## Shell profile commands

Because tab syncing is enabled by default, the shell choice below will stay in sync with the next shell tab group on this page.

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
fish_add_path $HOME/.local/bin
```
{{< /tab >}}
{{< /tabs >}}

## Entry check

{{< tabs >}}
{{< tab name="bash" selected=true >}}
```bash
hugo version
go version
```
{{< /tab >}}
{{< tab name="zsh" >}}
```zsh
hugo version
go version
```
{{< /tab >}}
{{< tab name="fish" >}}
```fish
hugo version
go version
```
{{< /tab >}}
{{< /tabs >}}

If the shell tabs stay aligned between the previous two sections, the synced-tab behavior is working as intended.
