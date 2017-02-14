import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { AuthGuard } from '../auth/auth.guard';
import { VolumeFormComponent } from './form/volume-form.component';
import { VolumeFormGroupComponent } from './form/volume-form-group.component';

const volumeRoutes: Routes = [{
    path: 'environments/:environment-uuid/versions/:version-uuid/volumes/:volume-uuid',
    component: VolumeFormGroupComponent,
    canActivate: [AuthGuard]
  },
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(volumeRoutes),
    SharedModule,
  ],
  declarations: [
    VolumeFormComponent,
    VolumeFormGroupComponent,
  ],
  exports: [
    VolumeFormGroupComponent,
    RouterModule,
  ],
})

export class VolumeModule {}
