import { Directive, TemplateRef, ViewContainerRef } from '@angular/core';

@Directive({
  selector: '[appAllowed]'
})
export class MockAllowedDirective {
  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
  ) {
    this.viewContainer.createEmbeddedView(this.templateRef);
  }
}
