import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';
import { VolumeModule } from '../volumes/volume.module';

import { AuthGuard } from '../auth/auth.guard';
import { EnvironmentDetailsComponent } from './details/environment-details.component';
import { VolumeFormComponent } from '../volumes/volume-form.component';

const environmentRoutes: Routes = [
  { path: 'environments/:uuid', component: EnvironmentDetailsComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(environmentRoutes),
    SharedModule,
    VolumeModule,
  ],
  declarations: [
    EnvironmentDetailsComponent,
    VolumeFormComponent,
  ],
  exports: [
    EnvironmentDetailsComponent,
    RouterModule,
  ],
})

export class EnvironmentModule {}
