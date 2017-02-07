import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ApplicationModule } from '../applications/application.module';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { DashboardComponent } from './dashboard.component';

const dashboardRoutes: Routes = [
  { path: '', component: DashboardComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    ApplicationModule,
    RouterModule.forChild(dashboardRoutes),
    SharedModule,
  ],
  declarations: [
    DashboardComponent,
  ],
  exports: [
    DashboardComponent,
    RouterModule,
  ],
})

export class DashboardModule {};
