import { NgModule }             from '@angular/core';
import { BrowserModule }        from '@angular/platform-browser';
import { FormsModule }          from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';

// App is our top level component
import { AppComponent }         from './app.component';
import { DashboardComponent }   from './dashboard/dashboard.component';
import { LoginComponent }       from './login/login.component';
import { MastheadComponent }    from './masthead/masthead.component';

const routes: Routes = [
    { path: 'dashboard', component: DashboardComponent },
    { path: 'login',     component: LoginComponent }
];

@NgModule({
    imports: [
        BrowserModule,
        FormsModule,
        RouterModule.forRoot(routes, { useHash: true })
    ],
    declarations: [
        AppComponent,
        DashboardComponent,
        LoginComponent,
        MastheadComponent
    ],
    bootstrap: [ AppComponent ]
})

export class AppModule {}
