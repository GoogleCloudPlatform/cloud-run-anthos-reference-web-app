import { ComponentFixture } from '@angular/core/testing';

export function SetFormValue(fixture: ComponentFixture<any>, selector: string, value: string) {
  const el = fixture.debugElement.nativeElement.querySelector(selector);
  el.value = value;
  el.dispatchEvent(new Event('input'));
}
