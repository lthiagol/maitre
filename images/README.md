# Images

**Every image the project uses lives here.** A screenshot, a banner, a diagram, a
render — if it is a picture, it is in this folder and nowhere else.

| | |
|---|---|
| [`logo/`](logo/) | the banner, the mark and the ASCII art |
| **Subfolders** | by subject. A new kind of image gets a new subfolder, not a new file at the root. |
| **Format** | **SVG** when the image is drawn — a banner, a diagram, a logo. It stays crisp, it stays diffable, and it is text. **PNG** when the image is *captured* — a screenshot, a photograph, a render of something that already exists. |
| **Where it is used** | the README embeds from here. The paths are relative to the repo root, and the folder moves with the file. |
| **Naming** | the subject, then the variant: `console-board.png`, `banner-plan-behind.svg`. No dates, no `final`, no `v2`. |

## Two rules for a screenshot

1. **Capture the tool, not a mock-up.** The console's own `--render` output is the
   truth; a drawing of it is a claim about it.
2. **Capture it wide.** The console is responsive, and a screenshot at 80 columns proves
   less than one at 120.

## The banner

The banner is [`logo/9-plan-behind.svg`](logo/9-plan-behind.svg). It is the card at the
top of the README, and it is the one picture the project is judged by first.
