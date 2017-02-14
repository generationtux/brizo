import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { EnvironmentDetailsComponent } from './details/environment-details.component';
import { ConfigurationComponent } from './configuration/configuration.component';

const environmentRoutes: Routes = [
  { path: 'environments/:uuid', component: EnvironmentDetailsComponent, canActivate: [AuthGuard] },
  { path: 'environments/:uuid/configuration', component: ConfigurationComponent, canActivate: [AuthGuard]},
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(environmentRoutes),
    SharedModule,
  ],
  declarations: [
    EnvironmentDetailsComponent,
    ConfigurationComponent,
  ],
  exports: [
    EnvironmentDetailsComponent,
    RouterModule,
  ],
})

export class EnvironmentModule {}
