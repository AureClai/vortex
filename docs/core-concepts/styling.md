# Styling in Vortex

Vortex features a comprehensive and type-safe "CSS-in-Go" styling system. This approach allows you to define your component's styles directly in Go, right alongside your component logic. This provides better modularity, reusability, and compile-time checks for your styles.

## Table of Contents

- [The `style` Package](#the-style-package)
- [Creating a Basic Style](#creating-a-basic-style)
- [Core Features](#core-features)
  - [Flexbox Layout](#flexbox-layout)
  - [Pseudo-classes](#pseudo-classes)
  - [Media Queries](#media-queries)
  - [Theming](#theming)
- [Benefits of CSS-in-Go](#benefits-of-css-in-go)
- [How It Works Under the Hood](#how-it-works-under-the-hood)
  - [1. Style Definition as Go Objects](#1-style-definition-as-go-objects)
  - [2. Pre-compilation and Caching](#2-pre-compilation-and-caching)
  - [3. Unique Class Name Generation](#3-unique-class-name-generation)
  - [4. On-Demand CSS Injection](#4-on-demand-css-injection)

## The `style` Package

All styling capabilities are provided by the `pkg/style` package. The core of the system is the `style.Style` struct, which you use to define a set of style rules.

## Creating a Basic Style

To style a component, you create a new style object using `style.New()`. This function accepts a series of style options using the functional options pattern.

Here is a simple example of styling a `Text` component:

```go
package main

import (
    "github.com/AureClai/vortex/pkg/component"
    "github.com/AureClai/vortex/pkg/style"
    "github.com/AureClai/vortex/pkg/vdom"
)

// Define a style for our title
var titleStyle = style.New(
    style.Color(style.ColorBlue),
    style.FontSize(style.Px(24)),
    style.FontWeight(style.FontWeightBold),
    style.Padding(style.PaddingAll, style.Px(10)),
)

// A simple component that uses the style
func App() *vdom.VNode {
    textComponent := component.NewText("Hello, Styled World!")
    textComponent.Style(titleStyle)
    return textComponent.Render()
}
```

In this example:
1.  We create `titleStyle` using `style.New()` with functional options like `style.Color()`, `style.FontSize()`, and `style.Padding()`.
2.  We create a text component using `component.NewText()`.
3.  We apply the style to the component using the `.Style()` method.
4.  We return the rendered component using `.Render()`.

Vortex automatically generates a unique class name for this style, injects the corresponding CSS into a `<style>` tag in the document's head, and applies the class to the component's HTML element.

## Core Features

The styling system is designed to be comprehensive and cover the modern CSS features you already know.

### Flexbox Layout

Vortex has first-class support for Flexbox. You can create complex and responsive layouts with ease.

```go
var containerStyle = style.New(
    style.Display(style.DisplayFlex),
    style.FlexDirection(style.FlexDirectionColumn),
    style.JustifyContent(style.JustifyContentCenter),
    style.AlignItems(style.AlignItemsCenter),
    style.WidthPx(300),
    style.HeightPx(200),
    style.BackgroundColor("#f0f0f0"),
)
```

### Pseudo-classes

You can define styles for pseudo-classes like `:hover`, `:focus`, or `:active` using functions like `OnHover()`, `OnFocus()`, and `OnActive()`.

```go
var buttonStyle = style.New(
    style.BackgroundColor("blue"),
    style.Color("white"),
    style.Padding(style.PaddingAll, style.Px(20)),
    style.OnHover(
        style.BackgroundColor("darkblue"),
    ),
)
```
This example will change the button's background color when the user hovers over it.

### Media Queries

Responsive design is a key part of modern web development. You can apply styles based on screen size using the `MediaQuery()` function.

```go
var responsiveContainer = style.New(
    style.WidthPx(600), // Default width
    style.MediaQuery(
        style.MediaQueryTypeMaxWidth, 
        "768px",
        style.WidthPx(300), // Width on smaller screens
    ),
)
```

### Theming

Vortex includes a theming system that allows you to define a consistent design language for your application. You can define colors, fonts, and spacing in a central theme and then reference them in your component styles.

This feature is more advanced and will be covered in its own guide.

## Benefits of CSS-in-Go

-   **Type Safety**: No more typos in property names or values. The Go compiler catches errors for you.
-   **Scoped Styles**: Styles are automatically scoped to the components they are applied to, preventing unintended side effects.
-   **Dynamic Styles**: Since styles are just Go code, you can easily create dynamic styles based on your component's state or props.
-   **Co-location**: Keeping your styles with your component logic makes your codebase easier to understand and maintain.

## How It Works Under the Hood

The CSS-in-Go system in Vortex is designed for performance and efficiency. It achieves this through a process of pre-compilation, unique class generation, and on-demand CSS injection.

### 1. Style Definition as Go Objects

When you define a style using `style.New(style.Color("blue"))`, you are not writing a string of CSS. Instead, you are constructing a Go `style.Style` object in memory. This object is a structured representation of your CSS rules.

This object-based representation is what allows for the type-safety and composability of the system.

### 2. Pre-compilation and Caching

Vortex includes a `PrecompilationEngine` that can be used to process these style objects ahead of time. When a style is "pre-compiled," the engine:
1.  Generates the final CSS string from the `style.Style` object.
2.  Generates a unique class name for that specific set of rules.
3.  Caches both the class name and the CSS string.

This means that for any given style, the CSS is generated only once.

### 3. Unique Class Name Generation

To ensure that styles don't conflict with each other, Vortex generates a unique class name for each distinct `style.Style` object. This name is a hash of the style's properties, resulting in a format like `vtx-a1b2c3d4`.

For example, `style.New(style.Color("blue"))` will always generate the same class name, while `style.New(style.Color("red"))` will generate a different one. This process ensures that styles are perfectly scoped and reusable.

### 4. On-Demand CSS Injection

Vortex is smart about how it adds CSS to the page. When the application starts, Vortex creates a single `<style id="vortex-styles">` tag in the `<head>` of your `index.html`.

The first time a component with a specific style (e.g., `.vtx-a1b2c3d4`) is rendered, the Vortex `Renderer` checks if the CSS for that class has already been injected.

-   If it has **not** been injected, the renderer appends the CSS rules to the content of the `<style>` tag.
-   If it **has** been injected, the renderer does nothing, avoiding duplication.

This on-demand injection mechanism, managed by the renderer's `updateStyle` and `processStyle` functions, ensures that only the CSS that is actually needed by the components on the screen is present in the DOM, keeping the stylesheet lean and performant.

This entire process is automatic and transparent to the developer. You simply create and apply style objects, and Vortex handles the efficient generation and management of the CSS behind the scenes.