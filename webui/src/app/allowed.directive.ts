import { Directive, Input, TemplateRef, ViewContainerRef } from '@angular/core';
import { LoginGuard } from './login.guard';

@Directive({
  selector: '[appAllowed]'
})
export class AllowedDirective {
  private hasView = false;

  constructor(
    private templateRef: TemplateRef<any>,
    private viewContainer: ViewContainerRef,
    private loginGuard: LoginGuard,
  ) { }

  @Input() set appAllowed(roles: Array<string>) {
    let allowed = false;
    if (roles.indexOf(this.loginGuard.userRole) >= 0) {
      allowed = true;
    }
    if (allowed && !this.hasView) {
      this.viewContainer.createEmbeddedView(this.templateRef);
      this.hasView = true;
    } else if (!allowed && this.hasView) {
      this.viewContainer.clear();
      this.hasView = false;
    }
  }

}
