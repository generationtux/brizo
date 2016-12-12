import { NgModule }                     from '@angular/core';
import { BrowserModule }                from '@angular/platform-browser';
import { FormsModule }                  from '@angular/forms';
import { RouterModule, Routes }         from '@angular/router';
import { HttpModule }                   from '@angular/http';
import { AUTH_PROVIDERS }               from 'angular2-jwt';

// App is our top level component
import { AppComponent }                 from './app.component';
import { DashboardComponent }           from './dashboard/dashboard.component';
import { LoginComponent }               from './login/login.component';
import { MastheadComponent }            from './masthead/masthead.component';
import { AuthService }                  from './services/auth.service';
import { AuthGuard }                    from './services/auth-guard.service';

const routes: Routes = [
    { path: '',          component: DashboardComponent, canActivate: [AuthGuard] },
    { path: 'dashboard', component: DashboardComponent, canActivate: [AuthGuard] },
    { path: 'login',     component: LoginComponent }
];

@NgModule({
    imports: [
        BrowserModule,
        FormsModule,
        HttpModule,
        RouterModule.forRoot(routes)
    ],
    declarations: [
        AppComponent,
        DashboardComponent,
        LoginComponent,
        MastheadComponent
    ],
    providers: [
        AUTH_PROVIDERS,
        AuthService,
        AuthGuard
    ],
    bootstrap: [ AppComponent ]
})

export class AppModule {}
