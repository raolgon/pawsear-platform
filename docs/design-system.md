# Pawsear Design System

## Direction

Pawsear should feel like a calm mobile-first operations tool. The interface is used while coordinating real pet care work, often during daytime and while moving between services, so clarity in bright environments is more important than decoration.

Design principles:

- Mobile-first by default.
- Light mode first.
- Cards are the main operational surface.
- Pet names and urgent care states must be immediately visible.
- Calendar context supports the workflow, but daily work queues lead.
- Desktop expands visibility; it should not be required for core tasks.
- Use Skeleton UI for component primitives and Tailwind v4 tokens for Pawsear-specific styling.

## Approved Visual Direction

The current UI direction is represented by these concept images:

- [`dashboard-v1.png`](design-concepts/dashboard-v1.png)
- [`household-mobile-v1.png`](design-concepts/household-mobile-v1.png)
- [`payments-v1.png`](design-concepts/payments-v1.png)

Shared shell:

- Desktop uses a persistent left navigation rail and a wide operational workspace.
- Mobile uses a compact top brand bar and fixed bottom navigation.
- Primary navigation is limited to Today, Calendar, Households, and Payments.
- Warm ivory is the application background; white is reserved for actionable surfaces.
- Deep navy carries hierarchy and primary actions. Muted gold identifies the brand and routine work.
- Coral is reserved for overdue or attention states; sage identifies completed work and received money.
- Corners are softly rounded, but cards stay compact and operational rather than decorative.

## Theme Strategy

Pawsear is light-mode first.

Reasons:

- The app will likely be used most during the day.
- Staff may check it outdoors, in building entrances, on walks, or in bright homes.
- Operational text, times, pet names, and urgent states must stay readable at a glance.
- Light surfaces make dense cards and lists easier to scan for this use case.

Dark mode is deferred until the core workflows are stable.

MVP theme scope:

- Build and QA light mode first.
- Do not design dark-only surfaces.
- Do not rely on dark backgrounds for brand personality.
- Keep enough contrast for outdoor/daytime use.
- Add dark mode later as a separate accessibility/usability pass.

## Tailwind v4 Color Tokens

The Pawsear palette is defined in:

`apps/web/src/routes/layout.css`

Tailwind class examples:

```text
bg-vanilla-custard-50
text-prussian-blue-950
border-light-gold-300
bg-tiger-flame-500
text-sandy-brown-800
```

## Palette Roles

### Foundation

Use light warm surfaces carefully so the app does not become visually heavy.

- Page background: `vanilla-custard-50`
- Raised surfaces/cards: white or `light-gold-50`
- Subtle dividers: `light-gold-200`
- Muted panels: `vanilla-custard-100`
- Main body text: `prussian-blue-950`
- Secondary text: neutral gray or `light-gold-800`

Avoid large dark panels in the MVP. Use dark colors mainly for text, icons, borders, and high-contrast actions.

### Primary Actions

Use Prussian Blue for primary navigation and high-confidence actions.

- Primary button: `prussian-blue-700`
- Primary button hover: `prussian-blue-800`
- Primary text on light surfaces: `prussian-blue-950`
- Focus ring: `prussian-blue-300`

### Attention And Urgency

Use Tiger Flame for urgent operational states. Do not overuse it.

- Overdue: `tiger-flame-600`
- Needs review: `sandy-brown-500`
- Destructive/cancel: `tiger-flame-700`
- Urgent background tint: `tiger-flame-50`

### Scheduling And Work States

Use the warm palette for routine operational status.

- Pending: `light-gold-500`
- Confirmed: `prussian-blue-600`
- In progress: `sandy-brown-500`
- Completed: `vanilla-custard-700`
- Cancelled/skipped: neutral gray from Tailwind default palette

### Pet And Care Cards

Pet cards should have strong readable text and restrained accents.

- Pet name: `prussian-blue-950`
- Household label: neutral gray or `light-gold-800`
- Care note accent: `sandy-brown-600`
- Medicine/overdue accent: `tiger-flame-600`

## Component Guidance

### Buttons

- Primary action: solid Prussian Blue.
- Secondary action: outlined or subtle light surface.
- Urgent action: Tiger Flame only when immediate attention is needed.
- Icon buttons should use lucide icons where possible.

### Cards

Cards are for operational items:

- booking card
- care task card
- boarding stay card
- message review card
- payment/charge card
- pet card

Card radius should stay modest: `rounded-lg` at most unless Skeleton defaults require otherwise.

Card hierarchy:

1. Pet or household subject.
2. Time/status/type.
3. Instruction or summary.
4. Actions.

### Status Labels

Status labels must include text, not color only.

Recommended mapping:

- Needs review: `bg-sandy-brown-100 text-sandy-brown-800`
- Overdue: `bg-tiger-flame-100 text-tiger-flame-800`
- Pending: `bg-light-gold-100 text-light-gold-800`
- Confirmed: `bg-prussian-blue-100 text-prussian-blue-800`
- Completed: `bg-vanilla-custard-200 text-vanilla-custard-900`

### Layout

Mobile:

- Use stacked cards.
- Keep primary action reachable.
- Use bottom sheets for focused creation/editing flows.
- Default to lists/queues over dense calendars.

Desktop:

- Use multi-column layouts for context.
- Keep the same mobile card model.
- Add filters, sidebars, and timelines as supporting surfaces.

## Accessibility

- Do not rely on color alone.
- Keep touch targets large enough for mobile use.
- Prefer visible labels for operational actions.
- Use high contrast for pet names, status labels, and due times.
- Keep urgent states visually distinct but not noisy.
- Test light mode in bright environments before investing in dark mode.

## Skeleton UI Usage

Use Skeleton for:

- Buttons.
- Inputs.
- Modals/sheets.
- Tabs.
- Menus.
- Toggles.
- Form controls.

Use Pawsear Tailwind tokens for:

- Brand color.
- Status labels.
- Card accents.
- Operational hierarchy.

Avoid making Skeleton theme choices before building the first static mockups. First validate the UI shape with these tokens, then formalize a Skeleton theme once patterns stabilize.
