import * as Yup from 'yup';

export const actionSchema = Yup.object().shape({
  feedback: Yup.string().required('Please fill out email')
});

export const lifecycleIdSchema = Yup.object().shape({
  lifecycleId: Yup.string().trim().nullable().length(7),
  expirationDateDay: Yup.string()
    .trim()
    .length(2)
    .required('Please include a day'),
  expirationDateMonth: Yup.string()
    .trim()
    .length(2)
    .required('Please include a month'),
  expirationDateYear: Yup.string()
    .trim()
    .length(4)
    .required('Please include a year'),
  scope: Yup.string().trim().required('Please include a scope'),
  nextSteps: Yup.string().trim(),
  feedback: Yup.string().trim().required('Please fill out email')
});

export const rejectIntakeSchema = Yup.object().shape({
  nextSteps: Yup.string().trim().required('Please include next steps'),
  reason: Yup.string().trim().required('Please include a reason'),
  feedback: Yup.string().trim().required('Please fill out email')
});

export const provideGRTFeedbackSchema = Yup.object().shape({
  grtFeedback: Yup.string()
    .trim()
    .required('Please include feedback for the business owner'),
  emailBody: Yup.string().trim().required('Please fill out email')
});
