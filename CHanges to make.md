Here is a production-grade system prompt designed to instruct an AI or developer to implement an admin-configurable Google Analytics engine that gracefully integrates with your application's public pages.

---

## System Prompt: Configurable Google Analytics Integration Module

### Objective

Design and implement an admin-configurable **Google Analytics Integration Module** within the administrative dashboard settings. This module must give administrators the ability to toggle Google Analytics tracking on or off globally across all public-facing pages, input their Google Measurement ID (Google Analytics 4), and ensure the tracking scripts are dynamically injected or entirely omitted based on the configuration state.

---

### 1. Database Configuration Schema

Create or extend a persistent configuration structure (such as a single-row `third_party_settings` table or global settings key-value store) to maintain the state:

* **`google_analytics_enabled`** (Boolean: `true` or `false`, default `false`)
* **`google_analytics_id`** (String, null allowed, strictly validated to match the GA4 format: `G-XXXXXXXXXX`)

---

### 2. Admin Settings Interface (UI/UX)

Provide a clear, dedicated configuration section within the administrative dashboard under "Third-Party Integrations":

* **Global Toggle:** A toggle switch or checkbox labeled **"Enable Google Analytics Tracking"**.
* **Conditional Inputs:**
* When the toggle is flipped to **Off**, the Measurement ID input field is disabled or hidden.
* When the toggle is flipped to **On**, display a text input field labeled **"Google Analytics Measurement ID (GA4)"** with a placeholder text showing `G-XXXXXXXXXX`.


* **Frontend Validation:** Enforce regex validation on the frontend to ensure the ID starts with `G-` followed by alphanumeric characters before allowing the administrator to click "Save Settings".

---

### 3. Frontend Script Injection Layout (Public Pages)

Develop a clean script rendering engine inside the root template layout file (e.g., in the main layout `<head>` block of your public-facing pages) that respects the configuration flag:

* **State Execution Logic:**
* **Case `false` (Disabled):** The backend must not output *any* Google Analytics tracking code to the browser. The page source must remain entirely free of GA scripts.
* **Case `true` (Enabled):** If the toggle is active and the `google_analytics_id` is populated, dynamically compile and inject the official Google Analytics async tag directly into the public `<head>` block:


```html
<!-- Google tag (gtag.js) -->
<script async src="https://www.googletagmanager.com/gtag/js?id={{GOOGLE_ANALYTICS_ID}}"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', '{{GOOGLE_ANALYTICS_ID}}');
</script>

```



---

### 4. Technical Architecture & Edge Cases

* **Server-Side Rendering (SSR) Optimization:** Ensure that the database call to check the analytics configuration state is heavily cached globally (e.g., using Redis or internal application memory caching). This avoids querying the database on every single public page visit just to check if the analytics tag should load.
* **Sanitization:** Strictly sanitize the `google_analytics_id` string on both the backend and frontend before rendering it to prevent any raw string manipulation or Cross-Site Scripting (XSS) injection vectors inside the public template layout.
* **Dynamic State Shifts:** If the administrator toggles the tracking engine from **On** to **Off**, the script must disappear on the very next page refresh without requiring a server reboot or complex manual cache clearing.