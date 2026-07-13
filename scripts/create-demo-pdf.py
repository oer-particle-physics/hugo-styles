#!/usr/bin/env python3

"""Create the deterministic one-page PDF used by the documentation showcase.

Requires ReportLab: `python3 -m pip install reportlab`.
"""

from __future__ import annotations

import argparse
from pathlib import Path

from reportlab.lib.colors import Color, HexColor, white
from reportlab.lib.pagesizes import A4
from reportlab.pdfbase.pdfmetrics import stringWidth
from reportlab.pdfgen import canvas


NAVY = HexColor("#0f172a")
BLUE = HexColor("#2563eb")
TEAL = HexColor("#0f766e")
SLATE = HexColor("#475569")
PALE_BLUE = HexColor("#eff6ff")
PALE_TEAL = HexColor("#f0fdfa")
LINE = HexColor("#cbd5e1")


def draw_wrapped_text(
    pdf: canvas.Canvas,
    text: str,
    x: float,
    y: float,
    width: float,
    *,
    font: str = "Helvetica",
    size: float = 9.5,
    leading: float = 13,
    color: Color = NAVY,
) -> float:
    words = text.split()
    lines: list[str] = []
    current = ""
    for word in words:
        candidate = f"{current} {word}".strip()
        if current and stringWidth(candidate, font, size) > width:
            lines.append(current)
            current = word
        else:
            current = candidate
    if current:
        lines.append(current)

    pdf.setFont(font, size)
    pdf.setFillColor(color)
    for line in lines:
        pdf.drawString(x, y, line)
        y -= leading
    return y


def draw_bullet(pdf: canvas.Canvas, text: str, x: float, y: float, width: float) -> float:
    pdf.setFillColor(TEAL)
    pdf.circle(x + 3, y + 3, 2.2, fill=1, stroke=0)
    return draw_wrapped_text(pdf, text, x + 13, y, width - 13, size=9.3, leading=12.5)


def build_pdf(output: Path) -> None:
    output.parent.mkdir(parents=True, exist_ok=True)
    width, height = A4
    pdf = canvas.Canvas(str(output), pagesize=A4, invariant=1)
    pdf.setTitle("Toy particle-analysis quick handout")
    pdf.setAuthor("Hugo Styles documentation")
    pdf.setSubject("A compact, deterministic example handout for the Hugo Styles feature demo")

    pdf.setFillColor(NAVY)
    pdf.rect(0, height - 118, width, 118, fill=1, stroke=0)
    pdf.setFillColor(HexColor("#5eead4"))
    pdf.setFont("Helvetica-Bold", 10)
    pdf.drawString(42, height - 38, "HUGO STYLES LIVE FEATURE DEMO")
    pdf.setFillColor(white)
    pdf.setFont("Helvetica-Bold", 25)
    pdf.drawString(42, height - 72, "Toy particle-analysis workflow")
    pdf.setFont("Helvetica", 11)
    pdf.setFillColor(HexColor("#cbd5e1"))
    pdf.drawString(42, height - 94, "A one-page handout for a reproducible counting exercise")

    margin = 42
    content_width = width - 2 * margin
    y = height - 148

    pdf.setFillColor(NAVY)
    pdf.setFont("Helvetica-Bold", 14)
    pdf.drawString(margin, y, "Learning target")
    y = draw_wrapped_text(
        pdf,
        "Apply a small event selection, compare signal and control counts, and record enough context for another learner to reproduce the result.",
        margin,
        y - 19,
        content_width,
        size=10.2,
        leading=14,
        color=SLATE,
    )

    y -= 10
    box_height = 82
    pdf.setFillColor(PALE_BLUE)
    pdf.roundRect(margin, y - box_height, content_width, box_height, 10, fill=1, stroke=0)
    steps = [
        ("1", "Load", "12 toy events"),
        ("2", "Select", "pT > 25 GeV"),
        ("3", "Count", "signal / control"),
        ("4", "Report", "cuts + result"),
    ]
    column_width = content_width / len(steps)
    for index, (number, label, detail) in enumerate(steps):
        left = margin + index * column_width
        center = left + column_width / 2
        pdf.setFillColor(BLUE if index < 3 else TEAL)
        pdf.circle(center, y - 23, 12, fill=1, stroke=0)
        pdf.setFillColor(white)
        pdf.setFont("Helvetica-Bold", 10)
        pdf.drawCentredString(center, y - 27, number)
        pdf.setFillColor(NAVY)
        pdf.setFont("Helvetica-Bold", 10)
        pdf.drawCentredString(center, y - 48, label)
        pdf.setFillColor(SLATE)
        pdf.setFont("Helvetica", 8.5)
        pdf.drawCentredString(center, y - 63, detail)
        if index < len(steps) - 1:
            pdf.setStrokeColor(LINE)
            pdf.setLineWidth(1.2)
            pdf.line(center + 18, y - 23, center + column_width - 18, y - 23)
    y -= box_height + 25

    gap = 18
    column = (content_width - gap) / 2
    left_x = margin
    right_x = margin + column + gap

    pdf.setFillColor(NAVY)
    pdf.setFont("Helvetica-Bold", 14)
    pdf.drawString(left_x, y, "Selection recipe")
    left_y = y - 21
    for bullet in [
        "Keep events with transverse momentum above 25 GeV.",
        "Use 80-100 GeV as the toy signal window.",
        "Use 60-80 GeV as the adjacent control window.",
    ]:
        left_y = draw_bullet(pdf, bullet, left_x, left_y, column)
        left_y -= 5

    pdf.setFillColor(PALE_TEAL)
    pdf.roundRect(right_x, y - 104, column, 104, 10, fill=1, stroke=0)
    pdf.setFillColor(TEAL)
    pdf.setFont("Helvetica-Bold", 10)
    pdf.drawString(right_x + 16, y - 21, "DETERMINISTIC EXAMPLE")
    pdf.setFillColor(NAVY)
    pdf.setFont("Helvetica-Bold", 23)
    pdf.drawString(right_x + 16, y - 51, "4 signal")
    pdf.setFont("Helvetica-Bold", 16)
    pdf.drawString(right_x + 16, y - 76, "2 control")
    pdf.setFillColor(SLATE)
    pdf.setFont("Helvetica", 8.6)
    pdf.drawString(right_x + 16, y - 93, "Matches the notebook output on the demo page.")

    y = min(left_y, y - 104) - 24
    pdf.setFillColor(NAVY)
    pdf.setFont("Helvetica-Bold", 14)
    pdf.drawString(margin, y, "Reproducibility check")
    y -= 21
    checks = [
        "Record the input sample and every numerical cut.",
        "Keep the code cell and its deterministic output together.",
        "State what the numbers demonstrate - and what they do not.",
    ]
    for index, check in enumerate(checks):
        row_y = y - index * 28
        pdf.setStrokeColor(BLUE)
        pdf.setLineWidth(1.3)
        pdf.roundRect(margin, row_y - 2, 11, 11, 2, fill=0, stroke=1)
        draw_wrapped_text(pdf, check, margin + 20, row_y, content_width - 20, size=9.4, leading=12.5)

    footer_y = 45
    pdf.setStrokeColor(LINE)
    pdf.line(margin, footer_y + 19, width - margin, footer_y + 19)
    pdf.setFillColor(SLATE)
    pdf.setFont("Helvetica", 8.5)
    pdf.drawString(margin, footer_y, "Local demo asset - builds and renders without network access")
    pdf.drawRightString(width - margin, footer_y, "CC BY 4.0 example")

    pdf.showPage()
    pdf.save()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "output",
        nargs="?",
        type=Path,
        default=Path("content/docs/hextra-features/particle-analysis-handout.pdf"),
    )
    args = parser.parse_args()
    build_pdf(args.output)


if __name__ == "__main__":
    main()
