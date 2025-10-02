# Animation in Vortex

Vortex features a multi-purposed animation engine. The core concepts of the animation engine is :

- It is called every requestFrame of the DOM to run at 60 fps.
- it modifies directly the DOM elements CSS inline properties

An animated compnent should implement the `AnimatedComponent` interface providing interface method `Update` as a `func(time.Duration)` type.

## The orchestrator

When creating a new AnimatedComponent the global animation orchestrator get a reference to it and can update the animation of the component when its own method `Update`is called.

## Animation types

Vortex provides 3 animation context

### Declarative CSS inline modification : The basics, simulation in hand

In the `Update` methode of an animated component, developper has access to the `props` and the `state`objects so that the animation at time t can depend on states or props passed to the element.

**Important** : the animation's logic parameters should be defined outside of the component's `state` to avoid reRendering the whole VDOM node at each `Update`call.
The `UpdateCSS(*style.Style)` methods is the type safe way to patch the style of an object.
For exemple, if the `background-color` property is animated, in Vortex, the usage is object.UpdateCSS(style.New().BackgroundColor(newcolor))

This animation context is the way to go for simulated behavior.

### Clips and animation state machine graphs : For animation that don't need complex logic to be implemented

Vortex provides a powerful state machine graph. A graph is composed of Nodes (with an animation Clip) and Edges that are transition. The graph has also parameters that enables to pass through transitions. Transitions conditions can also be triggered at the end of a Node Clip making.

#### Example : A button

While in basic CSS the buttons animation are defined with a compination of `transform`, `keyframes`, `animation`, `transition` properties, in Vortex the whole logic and animation can be done with an Animation graph.

- **Clips** : `"idle"` a small subtle idle animation loop, `"hoover"` a more impctful hoovered animation clip loop, `"just-pressed"` a very impactful animation with no looping, `"pressed"` an animation looping clip. Each clip has a set of tracks that represents a list of keyframes for the property animated
- **Transitions** (when we say a transition is triggered it correspond to the condition being fullfilled, if when a node change the following transition is already fullfilled, the node clip is not played) : `"idle"` to `"hover"` in based on the parameter `hovered bool` that is changed within the `On("mousein")` listener of the component. The reverse transition has juste the opposite condition. `"hover"` to `"just-pressed"` is triggered when `On("mousedown")`is triggered and the transition from `"just-pressed"` is triggered automatically at the end of `"just-pressed"` node clip. `On("mouseup")` trigger the transition from `"pressed"` to "`hover`". "`On("mouseout")`".

#### Other Example : Screen transition

The animation clip is based on the modification of the background color of div with zlevel above all. The idle state is nothing, the transition state has the clip that swipes the background gradient from left to right.

### SVG Animation : Fast and efficient morph your SVG components

Vortex embeds also SVG animation that provides some usefull tools to animate paths and more on SVG-based component.
