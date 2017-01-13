import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { HttpModule } from '@angular/http';

// App is our top level component
import { AppComponent } from './app.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { InviteComponent } from './invite/invite.component';
import { LoginComponent } from './login/login.component';
import { MastheadComponent } from './masthead/masthead.component';
import { AuthComponent } from './auth/auth.component';
import { AuthService } from './auth/auth.service';
import { AuthGuard } from './auth/auth.guard';
import { ApplicationComponent } from './resources/applications/application-details.component';
import { ApplicationCreationComponent } from './resources/applications/application-creation.component';
import { ApplicationListingComponent } from './resources/applications/application-listing.component';
import { EnvironmentService } from './resources/environments/environment.service';

const routes: Routes = [
    { path: '',                   component: DashboardComponent, canActivate: [AuthGuard] },
    { path: 'login',              component: LoginComponent },
    { path: 'auth',               component: AuthComponent },
    { path: 'applications/:uuid', component: ApplicationComponent, canActivate: [AuthGuard] },
    { path: 'users/invite',       component: InviteComponent, canActivate: [AuthGuard] },
];

@NgModule({
    imports: [
        BrowserModule,
        FormsModule,
        ReactiveFormsModule,
        RouterModule.forRoot(routes),
        HttpModule,
    ],
    providers: [
      AuthService,
      AuthGuard,
      EnvironmentService,
    ],
    declarations: [
        AuthComponent,
        AppComponent,
        ApplicationComponent,
        ApplicationCreationComponent,
        ApplicationListingComponent,
        DashboardComponent,
        InviteComponent,
        LoginComponent,
        MastheadComponent,
    ],
    bootstrap: [ AppComponent ]
})

export class AppModule {}
