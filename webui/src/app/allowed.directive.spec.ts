import { AllowedDirective } from './allowed.directive';
import { LoginGuard } from './login.guard';
import { TemplateRef, ViewContainerRef, Component } from '@angular/core';
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By } from '@angular/platform-browser';
import { AngularFireAuth } from '@angular/fire/auth';
import { RouterTestingModule } from '@angular/router/testing';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { of, Observable } from 'rxjs';
import { ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';

describe('AllowedDirective', () => {
  @Component({
    template: `<div *appAllowed="['admin']"><span data-testid="content">test</span></div>`
  })
  class TestComponent {
  }

  let loginGuard: LoginGuard;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AllowedDirective, TestComponent ],
      imports: [
        HttpClientTestingModule,
        RouterTestingModule,
        MatSnackBarModule,
      ],
      providers: [
        {
          provide: AngularFireAuth,
          useValue: {},
        }
      ]
    })
    .compileComponents();
    loginGuard = TestBed.inject(LoginGuard);
  }));

  it('should not display without allowed role', () => {
    loginGuard.userRole = '';
    const fixture = TestBed.createComponent(TestComponent);
    fixture.detectChanges();
    const directive = fixture.debugElement.queryAllNodes(By.directive(AllowedDirective));
    expect(directive).toBeTruthy();
    console.log(fixture.debugElement)
    const testComponent = fixture.debugElement.query(By.css('span[data-testid="content"]'));
    expect(testComponent).toBeFalsy();
  });

  it('should display with allowed role', () => {
    loginGuard.userRole = 'admin';
    const fixture = TestBed.createComponent(TestComponent);
    fixture.detectChanges();
    const directive = fixture.debugElement.queryAllNodes(By.directive(AllowedDirective));
    expect(directive).toBeTruthy();
    console.log(fixture.debugElement)
    const testComponent = fixture.debugElement.query(By.css('span[data-testid="content"]'));
    expect(testComponent).toBeTruthy();
  });


});
