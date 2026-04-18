const API = 'http://localhost:8080/api';

async function fetchJSON(path) {
  const res = await fetch(API + path);
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`);
  return res.json();
}

function renderList(id, items, label) {
  const ul = document.getElementById(id);
  ul.innerHTML = items.length
    ? items.map(i => `<li>${label(i)}</li>`).join('')
    : '<li><em>none</em></li>';
}

async function load() {
  try {
    const [users, items, transactions] = await Promise.all([
      fetchJSON('/users'),
      fetchJSON('/items'),
      fetchJSON('/transactions'),
    ]);
    renderList('users-list', users, u => `${u.username} — ${u.email}`);
    renderList('items-list', items, i => `${i.name} — $${i.price.toFixed(2)}`);
    renderList('transactions-list', transactions, t =>
      `TX#${t.id}: user ${t.user_id} × item ${t.item_id} (qty ${t.quantity})`
    );
  } catch (err) {
    console.error('Failed to load data:', err);
  }
}

load();
