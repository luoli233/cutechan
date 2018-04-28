/**
 * Follows lang.go structures.
 */

// tslint:disable-next-line:no-reference
/// <reference path="index.d.ts" />

import langs from "cc-langs";

interface LanguagePack {
  Plurals: { [key: string]: string[] };
  Forms: { [key: string]: string[2] };
  UI: { [key: string]: string };
}

// TODO(Kagami): Add support for per-user site language.
const siteLang: string = (window as any).config.defaultLang;
const pack: any = langs[siteLang];

/** Container of localization strings for current site language. Deprecated. */
export const ln: LanguagePack = { Plurals: pack.plurals, Forms: pack.forms, UI: pack.ui };

/** Gettext-alike helper. */
// TODO(Kagami): Rewrite ln.UI boilerplate to this.
export function _(s: string): string {
  return pack.ui[s] || s;
}

/** Printf-alike helper. */
export function printf(s: string, ...args: any[]): string {
  return s.replace(/%s/, () => args.shift());
}

export const months = [
  "Jan", "Feb", "Mar", "Apr", "May", "Jun",
  "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
];

export const days = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
