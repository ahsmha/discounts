import plotly.graph_objects as go
import plotly.io as pio

# Create flowchart using shapes and annotations
fig = go.Figure()

# Define colors for different stages
colors = {
    'input': '#1FB8CD',
    'process': '#FFC185', 
    'decision': '#ECEBD5',
    'output': '#5D878F'
}

# Add shapes for flowchart boxes
# Input box
fig.add_shape(type="rect", x0=0, y0=8, x1=4, y1=9, 
              fillcolor=colors['input'], line=dict(width=2))

# Process boxes
fig.add_shape(type="rect", x0=0, y0=6.5, x1=4, y1=7.5, 
              fillcolor=colors['process'], line=dict(width=2))

fig.add_shape(type="rect", x0=0, y0=5, x1=4, y1=6, 
              fillcolor=colors['decision'], line=dict(width=2))

fig.add_shape(type="rect", x0=0, y0=3.5, x1=4, y1=4.5, 
              fillcolor=colors['process'], line=dict(width=2))

# Output box
fig.add_shape(type="rect", x0=0, y0=2, x1=4, y1=3, 
              fillcolor=colors['output'], line=dict(width=2))

# Calculation detail boxes (right side)
fig.add_shape(type="rect", x0=5, y0=7, x1=9, y1=8, 
              fillcolor='#D2BA4C', line=dict(width=1))

fig.add_shape(type="rect", x0=5, y0=5.5, x1=9, y1=6.5, 
              fillcolor='#D2BA4C', line=dict(width=1))

fig.add_shape(type="rect", x0=5, y0=4, x1=9, y1=5, 
              fillcolor='#D2BA4C', line=dict(width=1))

# Add arrows
fig.add_shape(type="line", x0=2, y0=8, x1=2, y1=7.5, 
              line=dict(width=3, color="black"))
fig.add_shape(type="line", x0=2, y0=6.5, x1=2, y1=6, 
              line=dict(width=3, color="black"))
fig.add_shape(type="line", x0=2, y0=5, x1=2, y1=4.5, 
              line=dict(width=3, color="black"))
fig.add_shape(type="line", x0=2, y0=3.5, x1=2, y1=3, 
              line=dict(width=3, color="black"))

# Add text annotations
fig.add_annotation(x=2, y=8.5, text="Input: Cart + Customer<br>PUMA T-shirt x2 @ ₹1k<br>ICICI Bank Card", 
                   showarrow=False, font=dict(size=10))

fig.add_annotation(x=2, y=7, text="Retrieve Active Discounts<br>Brand: PUMA 40%<br>Category: T-shirt 10%<br>Bank: ICICI 10%", 
                   showarrow=False, font=dict(size=10))

fig.add_annotation(x=2, y=5.5, text="Filter & Validate<br>All 3 discounts apply", 
                   showarrow=False, font=dict(size=10))

fig.add_annotation(x=2, y=4, text="Apply Stacking Order<br>Brand → Category → Bank", 
                   showarrow=False, font=dict(size=10))

fig.add_annotation(x=2, y=2.5, text="Final Result<br>₹2000 → ₹972<br>51.4% savings", 
                   showarrow=False, font=dict(size=10))

# Calculation steps
fig.add_annotation(x=7, y=7.5, text="Step 1: Brand Discount<br>₹2000 - 40% = ₹1200", 
                   showarrow=False, font=dict(size=9))

fig.add_annotation(x=7, y=6, text="Step 2: Category Discount<br>₹1200 - 10% = ₹1080", 
                   showarrow=False, font=dict(size=9))

fig.add_annotation(x=7, y=4.5, text="Step 3: Bank Discount<br>₹1080 - 10% = ₹972", 
                   showarrow=False, font=dict(size=9))

# Set layout
fig.update_layout(
    title="Discount Calculation Process",
    xaxis=dict(range=[-0.5, 10], showgrid=False, showticklabels=False),
    yaxis=dict(range=[1, 10], showgrid=False, showticklabels=False),
    plot_bgcolor='white',
    showlegend=False
)

# Save the chart
fig.write_image("discount_calculation_flow.png")