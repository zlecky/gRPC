const out = document.getElementById("out");
const healthEl = document.getElementById("health");

function stamp() {
  return new Date().toLocaleTimeString();
}

async function call(method, path, body) {
  const started = performance.now();
  const opts = { method, headers: {} };
  if (body !== undefined) {
    opts.headers["Content-Type"] = "application/json";
    opts.body = JSON.stringify(body);
  }

  let res;
  let text;
  try {
    res = await fetch(path, opts);
    text = await res.text();
  } catch (err) {
    out.textContent = `[${stamp()}] ${method} ${path}\nnetwork error: ${err.message}`;
    return;
  }

  let pretty = text;
  try {
    pretty = JSON.stringify(JSON.parse(text), null, 2);
  } catch (_) {}

  const ms = Math.round(performance.now() - started);
  out.textContent =
    `[${stamp()}] ${method} ${path}\n` +
    `HTTP ${res.status} ${res.statusText} · ${ms}ms\n\n` +
    pretty;

  if (method === "POST" && res.ok) {
    try {
      const data = JSON.parse(text);
      const id = data?.user?.id;
      if (id) {
        document.getElementById("get-id").value = id;
        document.getElementById("update-id").value = id;
        document.getElementById("delete-id").value = id;
      }
    } catch (_) {}
  }
}

document.querySelector("[data-action=create]").addEventListener("click", () => {
  call("POST", "/api/users", {
    name: document.getElementById("create-name").value.trim(),
    email: document.getElementById("create-email").value.trim(),
  });
});

document.querySelector("[data-action=get]").addEventListener("click", () => {
  const id = document.getElementById("get-id").value.trim();
  call("GET", `/api/users/${encodeURIComponent(id)}`);
});

document.querySelector("[data-action=list]").addEventListener("click", () => {
  const size = document.getElementById("list-size").value || "10";
  const token = document.getElementById("list-token").value.trim();
  const q = new URLSearchParams({ page_size: size });
  if (token) q.set("page_token", token);
  call("GET", `/api/users?${q}`);
});

document.querySelector("[data-action=update]").addEventListener("click", () => {
  const id = document.getElementById("update-id").value.trim();
  call("PUT", `/api/users/${encodeURIComponent(id)}`, {
    name: document.getElementById("update-name").value.trim(),
    email: document.getElementById("update-email").value.trim(),
  });
});

document.querySelector("[data-action=delete]").addEventListener("click", () => {
  const id = document.getElementById("delete-id").value.trim();
  call("DELETE", `/api/users/${encodeURIComponent(id)}`);
});

document.getElementById("clear-out").addEventListener("click", () => {
  out.textContent = "等待请求…";
});

async function checkHealth() {
  try {
    const res = await fetch("/api/health");
    if (!res.ok) throw new Error(String(res.status));
    healthEl.textContent = "gateway ok";
    healthEl.classList.add("ok");
    healthEl.classList.remove("bad");
  } catch (_) {
    healthEl.textContent = "gateway down";
    healthEl.classList.add("bad");
    healthEl.classList.remove("ok");
  }
}

checkHealth();
setInterval(checkHealth, 8000);
