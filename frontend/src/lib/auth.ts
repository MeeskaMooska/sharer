const AUTH_STATUS_KEY = "sharer-auth-status";
const AUTH_EMAIL_KEY = "sharer-auth-email";
const AUTH_SCHOOL_KEY = "sharer-auth-school";

export type SupportedSchool = "nova" | "uva";

function hasWindow() {
  return typeof window !== "undefined";
}

export function isSignedIn() {
  if (!hasWindow()) {
    return false;
  }
  return window.localStorage.getItem(AUTH_STATUS_KEY) === "signed-in";
}

export function applySchoolTheme(school: SupportedSchool) {
  if (!hasWindow()) {
    return;
  }
  window.document.documentElement.setAttribute("data-school", school);
}

export function signInWithEmail(email: string, school: SupportedSchool) {
  if (!hasWindow()) {
    return;
  }
  window.localStorage.setItem(AUTH_STATUS_KEY, "signed-in");
  window.localStorage.setItem(AUTH_EMAIL_KEY, email.trim());
  window.localStorage.setItem(AUTH_SCHOOL_KEY, school);
  applySchoolTheme(school);
}

export function signOut() {
  if (!hasWindow()) {
    return;
  }
  window.localStorage.removeItem(AUTH_STATUS_KEY);
  window.localStorage.removeItem(AUTH_EMAIL_KEY);
  window.localStorage.removeItem(AUTH_SCHOOL_KEY);
  window.document.documentElement.removeAttribute("data-school");
}

export function getSignedInEmail() {
  if (!hasWindow()) {
    return "";
  }
  return window.localStorage.getItem(AUTH_EMAIL_KEY) ?? "";
}

export function getSignedInSchool(): SupportedSchool {
  if (!hasWindow()) {
    return "nova";
  }
  const school = window.localStorage.getItem(AUTH_SCHOOL_KEY);
  return school === "uva" ? "uva" : "nova";
}
