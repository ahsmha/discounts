import plotly.graph_objects as go
import plotly.express as px
import pandas as pd
import json

# Data from the provided JSON with exact colors
layers_data = [
    {"name": "Presentation", "description": "HTTP/gRPC API", "components": ["REST API", "gRPC API"], "y": 5, "color": "#4A90E2"},
    {"name": "Application", "description": "cmd/server", "components": ["main.go", "Bootstrap"], "y": 4, "color": "#7ED321"},
    {"name": "Business Logic", "description": "internal/service", "components": ["DiscountSvc", "Calculator", "Validator"], "y": 3, "color": "#F5A623"},
    {"name": "Data Access", "description": "internal/repo", "components": ["DiscountRepo", "MemoryRepo"], "y": 2, "color": "#D0021B"},
    {"name": "Infrastructure", "description": "pkg/errors", "components": ["Errors", "Ext Deps"], "y": 1, "color": "#9013FE"}
]

discount_types = [
    {"name": "Brand/Category", "step": 1, "priority": "Highest", "example": "PUMA 40% off"},
    {"name": "Voucher Codes", "step": 2, "priority": "Medium", "example": "SUPER69 69%"},
    {"name": "Bank Offers", "step": 3, "priority": "Lowest", "example": "ICICI 10%"}
]

solid_principles = [
    {"principle": "SRP", "desc": "Single Job"},
    {"principle": "OCP", "desc": "Extensible"},
    {"principle": "DIP", "desc": "Abstractions"}
]

# Create figure
fig = go.Figure()

# Add layer rectangles and components
for i, layer in enumerate(layers_data):
    color = layer["color"]
    
    # Add layer background rectangle
    fig.add_shape(
        type="rect",
        x0=-0.5, y0=layer["y"]-0.4,
        x1=6.5, y1=layer["y"]+0.4,
        line=dict(color=color, width=3),
        fillcolor=color,
        opacity=0.2,
        layer="below"
    )
    
    # Add layer name
    fig.add_trace(go.Scatter(
        x=[-0.2],
        y=[layer["y"]],
        mode='text',
        text=[layer["name"]],
        textposition="middle right",
        textfont=dict(size=12, color=color, family="Arial Black"),
        showlegend=False,
        cliponaxis=False
    ))
    
    # Add components as rectangles with text
    for j, component in enumerate(layer["components"]):
        fig.add_trace(go.Scatter(
            x=[j + 1.5],
            y=[layer["y"]],
            mode='markers+text',
            text=[component],
            textposition="middle center",
            textfont=dict(size=11, color="white", family="Arial"),
            marker=dict(
                size=90,
                color=color,
                line=dict(color="white", width=3),
                symbol="square",
                opacity=0.9
            ),
            showlegend=False,
            cliponaxis=False
        ))

# Add prominent data flow arrows (vertical flow)
arrow_props = dict(arrowhead=3, arrowsize=2, arrowwidth=3, arrowcolor="#13343B")

# Main processing flow arrows
for i in range(len(layers_data)-1):
    fig.add_annotation(
        x=3.5, y=layers_data[i]["y"]-0.6,
        text="",
        showarrow=True,
        ax=3.5, ay=layers_data[i]["y"]-0.2,
        **arrow_props
    )

# Input arrow
fig.add_annotation(
    x=1.5, y=5.6,
    text="INPUT",
    showarrow=True,
    arrowhead=3,
    arrowsize=2,
    arrowwidth=3,
    arrowcolor="#1FB8CD",
    ax=1.5, ay=6.2,
    font=dict(size=12, color="#1FB8CD", family="Arial Black")
)

# Output arrow
fig.add_annotation(
    x=5.5, y=0.4,
    text="OUTPUT",
    showarrow=True,
    arrowhead=3,
    arrowsize=2,
    arrowwidth=3,
    arrowcolor="#1FB8CD",
    ax=5.5, ay=-0.2,
    font=dict(size=12, color="#1FB8CD", family="Arial Black")
)

# Add discount stacking order section
discount_colors = ["#B4413C", "#964325", "#944454"]
for i, discount in enumerate(discount_types):
    y_pos = 4.5 - i * 1.2
    
    # Discount type box
    fig.add_trace(go.Scatter(
        x=[8.5],
        y=[y_pos],
        mode='markers+text',
        text=[f"Step {discount['step']}<br>{discount['name']}<br>{discount['example']}"],
        textposition="middle center",
        textfont=dict(size=10, color="white", family="Arial"),
        marker=dict(
            size=120,
            color=discount_colors[i],
            line=dict(color="white", width=2),
            symbol="square",
            opacity=0.9
        ),
        showlegend=False,
        cliponaxis=False
    ))
    
    # Priority indicator
    fig.add_trace(go.Scatter(
        x=[10],
        y=[y_pos],
        mode='text',
        text=[discount['priority']],
        textposition="middle center",
        textfont=dict(size=11, color=discount_colors[i], family="Arial Black"),
        showlegend=False,
        cliponaxis=False
    ))

# Add SOLID principles section
solid_colors = ["#5D878F", "#D2BA4C", "#ECEBD5"]
for i, principle in enumerate(solid_principles):
    fig.add_trace(go.Scatter(
        x=[8.5],
        y=[0.5 + i * 0.4],
        mode='markers+text',
        text=[f"{principle['principle']}<br>{principle['desc']}"],
        textposition="middle center",
        textfont=dict(size=9, color="black", family="Arial"),
        marker=dict(
            size=60,
            color=solid_colors[i],
            line=dict(color="white", width=2),
            symbol="circle",
            opacity=0.8
        ),
        showlegend=False,
        cliponaxis=False
    ))

# Add section headers
fig.add_trace(go.Scatter(
    x=[8.5],
    y=[5.5],
    mode='text',
    text=["Discount Stacking"],
    textposition="middle center",
    textfont=dict(size=14, color="#13343B", family="Arial Black"),
    showlegend=False,
    cliponaxis=False
))

fig.add_trace(go.Scatter(
    x=[8.5],
    y=[2.3],
    mode='text',
    text=["SOLID Principles"],
    textposition="middle center",
    textfont=dict(size=14, color="#13343B", family="Arial Black"),
    showlegend=False,
    cliponaxis=False
))

# Update layout
fig.update_layout(
    title="Unifize Discount Service Architecture",
    xaxis=dict(
        range=[-1, 11],
        showgrid=False,
        showticklabels=False,
        title=""
    ),
    yaxis=dict(
        range=[-0.5, 6.5],
        showgrid=False,
        showticklabels=False,
        title=""
    ),
    plot_bgcolor='rgba(0,0,0,0)',
    paper_bgcolor='rgba(0,0,0,0)'
)

# Save the chart
fig.write_image("unifize_architecture.png")
print("Updated chart saved successfully!")